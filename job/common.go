package job

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"go.uber.org/zap"
)

type env interface {
	Logger() *zap.Logger
	TaskService() service.TaskService
	TaskInputRequestTopic() topic.Topic[bus.TaskInputRequestEvent]
	TaskFinishTopic() topic.Topic[bus.TaskFinishEvent]
}
