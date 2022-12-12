package job

import (
	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/service"
	"go.uber.org/zap"
)

type transactionsCheckerJob struct {
	logger                *zap.Logger
	taskService           service.TaskService
	taskInputRequestTopic topic.Topic[bus.TaskInputRequestEvent]
}

func NewTransactionsCheckerJob(e env) *transactionsCheckerJob {
	return &transactionsCheckerJob{
		logger:                e.Logger(),
		taskService:           e.TaskService(),
		taskInputRequestTopic: e.TaskInputRequestTopic(),
	}
}

func (j *transactionsCheckerJob) ID() string         { return "transactions-checker" }
func (j *transactionsCheckerJob) CronConfig() string { return "*/5 * * * *" }

func (j *transactionsCheckerJob) Handler() {
	tasks := j.taskService.FindNotFinishedTasksThatDontHaveIncompletedTransactions()

	for _, task := range tasks {
		j.taskInputRequestTopic.Publish(bus.TaskInputRequestEvent{TaskId: task.Id})
	}
}
