package listener

import (
	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/service"
	"go.uber.org/zap"
)

type FuzzerListener interface {
	StartListening()
}

type fuzzerListener struct {
	cfg                   *config.Config
	logger                *zap.Logger
	taskInputRequestTopic topic.Topic[bus.TaskInputRequestEvent]
	taskService           service.TaskService
}

func NewFuzzerListener(e env) *fuzzerListener {
	return &fuzzerListener{
		cfg:                   e.Config(),
		logger:                e.Logger(),
		taskInputRequestTopic: e.TaskInputRequestTopic(),
		taskService:           e.TaskService(),
	}
}

func (l *fuzzerListener) StartListening() {
	l.taskInputRequestTopic.Subscribe(l.processEvent)
}

func (l *fuzzerListener) processEvent(evt bus.TaskInputRequestEvent) {
	task, err := l.taskService.Get(evt.TaskId)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving task: %v", err)
		return
	}

	if task.Status != common.TASK_RUNNING {
		l.logger.Sugar().Infof("the task %s is not running", task.Id)
		return
	}

}
