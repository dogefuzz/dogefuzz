package job

import (
	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/service"
	"go.uber.org/zap"
)

type env interface {
	Logger() *zap.Logger
	TaskService() service.TaskService
	TaskInputRequestTopic() topic.Topic[bus.TaskInputRequestEvent]
	TaskFinishTopic() topic.Topic[bus.TaskFinishEvent]
}
