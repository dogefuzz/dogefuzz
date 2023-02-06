package listener

import (
	"context"
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

	transactions, err := l.transactionService.FindByTaskId(task.Id)
	if err != nil {
		l.logger.Sugar().Errorf("an error occurred when retrieving transactions of this task: %v", err)
		return
	}

	transactionsReports := make([]common.TransactionReport, len(transactions))
	aggregatedWeakneses := make([]string, 0)
	for idx, transaction := range transactions {
		transactionsReports[idx] = l.buildTransactionReport(transaction)
		aggregatedWeakneses = append(aggregatedWeakneses, transaction.DetectedWeaknesses...)
	}

	report := common.TaskReport{
		TimeElapsed:        task.Expiration.Sub(task.StartTime),
		ContractName:       contract.Name,
		Coverage:           coverage.ComputeCoverage(contract.CFG, task.AggregatedExecutedInstructions),
		CoverageByTime:     l.computeCoverageByTime(contract.CFG, transactions),
		MinDistance:        distance.ComputeMinDistance(contract.DistanceMap, task.AggregatedExecutedInstructions),
		MinDistanceByTime:  l.computeMinDistanceByTime(contract.DistanceMap, transactions),
		Transactions:       transactionsReports,
		DetectedWeaknesses: common.GetUniqueSlice(aggregatedWeakneses),
	}
	err = l.reporterService.SendReport(ctx, report)
	if err != nil {
		l.logger.Sugar().Errorf("the report could not been sent: %v", err)
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

func (l *reporterListener) buildTransactionReport(transaction *dto.TransactionDTO) common.TransactionReport {
	return common.TransactionReport{
		Timestamp:            transaction.Timestamp,
		BlockchainHash:       transaction.BlockchainHash,
		Inputs:               transaction.Inputs,
		DetectedWeaknesses:   transaction.DetectedWeaknesses,
		ExecutedInstructions: transaction.ExecutedInstructions,
		DeltaCoverage:        transaction.DeltaCoverage,
		DeltaMinDistance:     transaction.DeltaMinDistance,
	}
}
