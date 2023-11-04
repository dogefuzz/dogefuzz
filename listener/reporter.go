package listener

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/coverage"
	"github.com/dogefuzz/dogefuzz/pkg/distance"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type reporterListener struct {
	logger             *zap.Logger
	taskFinishTopic    interfaces.Topic[bus.TaskFinishEvent]
	transactionService interfaces.TransactionService
	taskService        interfaces.TaskService
	functionService    interfaces.FunctionService
	contractService    interfaces.ContractService
	reporterService    interfaces.ReporterService
}

func NewReporterListener(e Env) *reporterListener {
	return &reporterListener{
		logger:             e.Logger(),
		taskFinishTopic:    e.TaskFinishTopic(),
		transactionService: e.TransactionService(),
		taskService:        e.TaskService(),
		functionService:    e.FunctionService(),
		contractService:    e.ContractService(),
		reporterService:    e.ReporterService(),
	}
}

func (l *reporterListener) Name() string { return "reporter" }
func (l *reporterListener) StartListening(ctx context.Context) {
	handler := func(evt bus.TaskFinishEvent) { l.processEvent(ctx, evt) }
	l.taskFinishTopic.Subscribe(handler)
	<-ctx.Done()
	l.taskFinishTopic.Unsubscribe(handler)
}

func (l *reporterListener) processEvent(ctx context.Context, evt bus.TaskFinishEvent) {
	l.logger.Debug("processing TaskFinishEvent...")

	task, err := l.taskService.Get(evt.TaskId)
	if err != nil {
		l.logger.Sugar().Errorf("task could not be retrieved: %v", err)
		return
	}

	contract, err := l.contractService.FindByTaskId(task.Id)
	if err != nil {
		l.logger.Sugar().Errorf("contract could not be retrieved: %v", err)
		return
	}
	if contract == nil {
		l.logger.Sugar().Errorf("no contract was found: %v", err)
		return
	}

	transactions, err := l.transactionService.FindDoneByTaskId(task.Id)
	if err != nil {
		l.logger.Sugar().Errorf("an error occurred when retrieving transactions of this task: %v", err)
		return
	}

	heatmap := l.initHeatMap(contract.CFG)
	aggregatedWeakneses := make([]string, 0)
	var totalCoverage float64 = 0
	var averageCoverage float64 = 0
	var criticalInstructionsHits uint64 = 0

	for _, transaction := range transactions {
		l.updateHeatMapWithExecutedInstructions(heatmap, transaction.ExecutedInstructions)
		aggregatedWeakneses = append(aggregatedWeakneses, transaction.DetectedWeaknesses...)
		criticalInstructionsHits += transaction.CriticalInstructionsHits
		totalCoverage += float64(transaction.Coverage)
	}

	if len(transactions) == 0 || math.IsNaN(totalCoverage) {
		averageCoverage = 0.0
	} else {
		averageCoverage = totalCoverage / float64(len(transactions))
	}

	report := common.TaskReport{
		TaskId:                   task.Id,
		TaskStatus:               task.Status,
		TimeElapsed:              task.Expiration.Sub(task.StartTime),
		ContractName:             contract.Name,
		TotalInstructions:        uint64(len(contract.CFG.GetEdgesPCs())),
		Coverage:                 coverage.ComputeCoverage(contract.CFG, task.AggregatedExecutedInstructions),
		CoverageByTime:           l.computeCoverageByTime(contract.CFG, transactions),
		MinDistance:              distance.ComputeMinDistance(contract.DistanceMap, task.AggregatedExecutedInstructions),
		MinDistanceByTime:        l.computeMinDistanceByTime(contract.DistanceMap, transactions),
		DetectedWeaknesses:       common.GetUniqueSlice(aggregatedWeakneses),
		CriticalInstructionsHits: criticalInstructionsHits,
		AverageCoverage:          averageCoverage,
		Instructions:             l.buildInstructionsMap(contract.CFG),
		InstructionHitsHeatMap:   heatmap,
	}
	err = l.reporterService.SendReport(ctx, report)
	if err != nil {
		l.logger.Sugar().Errorf("the report could not been sent: %v", err)
		return
	}

	if task.Status != common.TASK_DEPLOY_ERROR {
		task.Status = common.TASK_DONE
	}

	err = l.taskService.Update(task)
	if err != nil {
		l.logger.Sugar().Errorf("the task %s could not be updated", task.Id)
		return
	}

	l.logger.Sugar().Infof("the execution report for the task %s was successfully sent", task.Id)
}

func (l *reporterListener) computeCoverageByTime(cfg common.CFG, transactions []*dto.TransactionDTO) common.TimeSeriesData {
	return l.computeTimeseriesFromTransactions(transactions, func(aggregatedInstructions []string) uint64 {
		return coverage.ComputeCoverage(cfg, aggregatedInstructions)
	})
}

func (l *reporterListener) computeMinDistanceByTime(distanceMap common.DistanceMap, transactions []*dto.TransactionDTO) common.TimeSeriesData {
	return l.computeTimeseriesFromTransactions(transactions, func(aggregatedInstructions []string) uint64 {
		return distance.ComputeMinDistance(distanceMap, aggregatedInstructions)
	})
}

func (l *reporterListener) computeTimeseriesFromTransactions(transactions []*dto.TransactionDTO, computeHandler func([]string) uint64) common.TimeSeriesData {
	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].Timestamp.Before(transactions[j].Timestamp)
	})

	result := common.TimeSeriesData{
		X: make([]time.Time, 0),
		Y: make([]uint64, 0),
	}
	aggregatedExecutedInstructions := make([]string, 0)
	for _, tx := range transactions {
		result.X = append(result.X, tx.Timestamp)
		aggregatedExecutedInstructions = common.MergeSortedSlices(aggregatedExecutedInstructions, tx.ExecutedInstructions)
		result.Y = append(result.Y, computeHandler(aggregatedExecutedInstructions))
	}

	return result
}

func (l *reporterListener) buildInstructionsMap(cfg common.CFG) map[string]string {
	instructionsMap := make(map[string]string)
	for _, block := range cfg.Blocks {
		for pc, instruction := range block.Instructions {
			instructionsMap[pc] = instruction
		}
	}
	return instructionsMap
}

func (l *reporterListener) initHeatMap(cfg common.CFG) map[string]uint64 {
	heatmap := make(map[string]uint64)
	for _, block := range cfg.Blocks {
		for pc, _ := range block.Instructions {
			heatmap[pc] = 0
		}
	}
	return heatmap
}

func (l *reporterListener) updateHeatMapWithExecutedInstructions(heatmap map[string]uint64, executedInstructions []string) map[string]uint64 {
	for _, pc := range executedInstructions {
		if _, ok := heatmap[pc]; !ok {
			continue
		}
		heatmap[pc]++
	}
	return heatmap
}
