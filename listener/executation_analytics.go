package listener

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/coverage"
	"github.com/dogefuzz/dogefuzz/pkg/distance"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"go.uber.org/zap"
)

type executionAnalyticsListener struct {
	logger                   *zap.Logger
	instrumentExecutionTopic topic.Topic[bus.InstrumentExecutionEvent]
	contractService          service.ContractService
	transactionService       service.TransactionService
	taskService              service.TaskService
}

func NewExecutionAnalyticsListener(e env) *executionAnalyticsListener {
	return &executionAnalyticsListener{
		logger:                   e.Logger(),
		instrumentExecutionTopic: e.InstrumentExecutionTopic(),
		contractService:          e.ContractService(),
		transactionService:       e.TransactionService(),
		taskService:              e.TaskService(),
	}
}

func (l *executionAnalyticsListener) StartListening() {
	l.instrumentExecutionTopic.Subscribe(l.processEvent)
}

func (l *executionAnalyticsListener) processEvent(evt bus.InstrumentExecutionEvent) {
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

	constract, err := l.contractService.Get(task.ContractId)
	if err != nil {
		l.logger.Sugar().Errorf("contract could not be retrieved: %v", err)
		return
	}

	transaction.DeltaCoverage = coverage.ComputeDeltaCoverage(constract.CFG, transaction.ExecutedInstructions, task.AggregatedExecutedInstructions)
	transaction.DeltaMinDistance = distance.ComputeDeltaMinDistance(constract.DistanceMap, transaction.ExecutedInstructions, task.AggregatedExecutedInstructions)
	transaction.Status = common.TRANSACTION_DONE
	err = l.transactionService.Update(transaction)
	if err != nil {
		l.logger.Sugar().Errorf("transaction could not be saved: %v", err)
		return
	}

	task.AggregatedExecutedInstructions = common.MergeSortedSlices(transaction.ExecutedInstructions, task.AggregatedExecutedInstructions)
	err = l.taskService.Update(task)
	if err != nil {
		l.logger.Sugar().Errorf("task could not be saved: %v", err)
		return
	}
}
