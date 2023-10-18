package job

import (
	"time"

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

	maxRetries := 5
	var find_error error
	for retries := 0; retries < maxRetries; retries++ {
		tasks, find_error := j.taskService.FindNotFinishedThatHaveDeployedContractAndLimitedPendingTransactions(j.cfg.FuzzerConfig.PendingTransactionsThreshold)
		if find_error != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		for _, task := range tasks {
			j.taskInputRequestTopic.Publish(bus.TaskInputRequestEvent{TaskId: task.Id})
		}

		return
	}
	j.logger.Sugar().Errorf("an error occured when retrieving tasks that are still running: %v", find_error)
}
