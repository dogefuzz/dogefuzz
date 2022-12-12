package job

import (
	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/service"
	"go.uber.org/zap"
)

type tasksCheckerJob struct {
	logger          *zap.Logger
	taskService     service.TaskService
	taskFinishTopic topic.Topic[bus.TaskFinishEvent]
}

func NewTasksCheckerJob(e env) *tasksCheckerJob {
	return &tasksCheckerJob{
		logger:          e.Logger(),
		taskService:     e.TaskService(),
		taskFinishTopic: e.TaskFinishTopic(),
	}
}

func (j *tasksCheckerJob) ID() string         { return "tasks-checker" }
func (j *tasksCheckerJob) CronConfig() string { return "*/5 * * * *" }

func (j *tasksCheckerJob) Handler() {
	tasks := j.taskService.FindNotFinishedAndExpired()

	for _, task := range tasks {
		j.taskFinishTopic.Publish(bus.TaskFinishEvent{TaskId: task.Id})
	}
}
