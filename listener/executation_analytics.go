package listener

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/coverage"
	"github.com/dogefuzz/dogefuzz/pkg/distance"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type executionAnalyticsListener struct {
	logger                   *zap.Logger
	instrumentExecutionTopic interfaces.Topic[bus.InstrumentExecutionEvent]
	contractService          interfaces.ContractService
	transactionService       interfaces.TransactionService
	taskService              interfaces.TaskService
}

func NewExecutionAnalyticsListener(e Env) *executionAnalyticsListener {
	return &executionAnalyticsListener{
		logger:                   e.Logger(),
		instrumentExecutionTopic: e.InstrumentExecutionTopic(),
		contractService:          e.ContractService(),
		transactionService:       e.TransactionService(),
		taskService:              e.TaskService(),
	}
}

func (l *executionAnalyticsListener) Name() string { return "execution-analytics" }
func (l *executionAnalyticsListener) StartListening(ctx context.Context) {
	handler := func(evt bus.InstrumentExecutionEvent) { l.processEvent(ctx, evt) }
	l.instrumentExecutionTopic.Subscribe(handler)
	<-ctx.Done()
	l.instrumentExecutionTopic.Unsubscribe(handler)
}

func (l *executionAnalyticsListener) processEvent(ctx context.Context, evt bus.InstrumentExecutionEvent) {
	l.logger.Debug("processing InstrumentExecutionEvent...")

	transaction, err := l.transactionService.Get(evt.TransactionId)
	if err != nil {
		l.logger.Sugar().Errorf("transaction could not be retrieved: %v", err)
		return
	}

	task, err := l.taskService.Get(transaction.TaskId)
	if err != nil {
		l.logger.Sugar().Errorf("task could not be retrieved: %v", err)
		return
	}

	contract, err := l.contractService.FindByTaskId(task.Id)
	if err != nil {
		l.logger.Sugar().Errorf("contract could not be retrieved: %v", err)
		return
	}

	transaction.DeltaCoverage = coverage.ComputeDeltaCoverage(contract.CFG, transaction.ExecutedInstructions, task.AggregatedExecutedInstructions)
	transaction.DeltaMinDistance = distance.ComputeDeltaMinDistance(contract.DistanceMap, transaction.ExecutedInstructions, task.AggregatedExecutedInstructions)
	transaction.Status = common.TRANSACTION_DONE
	err = l.transactionService.Update(transaction)
	if err != nil {
		l.logger.Sugar().Errorf("transaction could not be saved: %v", err)
		return
	}
	l.logger.Sugar().Debugf("transaction %s has achieved additional %d%% coverage", transaction.Id, uint64(100*float64(transaction.DeltaCoverage)/float64(len(contract.CFG.GetEdgesPCs()))))
	l.logger.Sugar().Debugf("transaction %s has reduce distance in %d", transaction.Id, transaction.DeltaMinDistance)

	task.AggregatedExecutedInstructions = common.MergeSortedSlices(transaction.ExecutedInstructions, task.AggregatedExecutedInstructions)
	err = l.taskService.Update(task)
	if err != nil {
		l.logger.Sugar().Errorf("task could not be saved: %v", err)
		return
	}
}
