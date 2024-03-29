package job

import (
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type transactionsCheckerJob struct {
	logger                *zap.Logger
	cfg                   *config.Config
	taskService           interfaces.TaskService
	transactionService    interfaces.TransactionService
	taskInputRequestTopic interfaces.Topic[bus.TaskInputRequestEvent]
}

func NewTransactionsCheckerJob(e Env) *transactionsCheckerJob {
	return &transactionsCheckerJob{
		logger:                e.Logger(),
		cfg:                   e.Config(),
		taskService:           e.TaskService(),
		transactionService:    e.TransactionService(),
		taskInputRequestTopic: e.TaskInputRequestTopic(),
	}
}

func (j *transactionsCheckerJob) Id() string         { return "transactions-checker" }
func (j *transactionsCheckerJob) CronConfig() string { return "* * * * * *" }

func (j *transactionsCheckerJob) Handler() {
	tasks, err := j.taskService.FindNotFinishedThatHaveDeployedContractAndLimitedPendingTransactions(j.cfg.FuzzerConfig.PendingTransactionsThreshold)
	if err != nil {
		j.logger.Sugar().Errorf("an error occured when retrieving tasks that are still running: %v", err)
		return
	}

	for _, task := range tasks {
		j.taskInputRequestTopic.Publish(bus.TaskInputRequestEvent{TaskId: task.Id})
	}
}
