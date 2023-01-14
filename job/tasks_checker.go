package job

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type tasksCheckerJob struct {
	logger          *zap.Logger
	taskService     interfaces.TaskService
	taskFinishTopic interfaces.Topic[bus.TaskFinishEvent]
}

func NewTasksCheckerJob(e Env) *tasksCheckerJob {
	return &tasksCheckerJob{
		logger:          e.Logger(),
		taskService:     e.TaskService(),
		taskFinishTopic: e.TaskFinishTopic(),
	}
}

func (j *tasksCheckerJob) Id() string         { return "tasks-checker" }
func (j *tasksCheckerJob) CronConfig() string { return "*/5 * * * *" }

func (j *tasksCheckerJob) Handler() {
	tasks, err := j.taskService.FindNotFinishedAndExpired()
	if err != nil {
		j.logger.Sugar().Errorf("an error occured when retrieving tasks to be finished: %v", err)
		return
	}

	for _, task := range tasks {
		j.taskFinishTopic.Publish(bus.TaskFinishEvent{TaskId: task.Id})
	}
}
