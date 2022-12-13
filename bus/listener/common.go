package listener

import (
	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/service"
	"go.uber.org/zap"
)

type env interface {
	Config() *config.Config
	Logger() *zap.Logger
	TaskStartTopic() topic.Topic[bus.TaskStartEvent]
	TaskInputRequestTopic() topic.Topic[bus.TaskInputRequestEvent]
	TaskService() service.TaskService
	GethService() service.GethService
	VandalService() service.VandalService
	ContractService() service.ContractService
}
