package listener

import (
	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/service"
	"go.uber.org/zap"
)

type env interface {
	Logger() *zap.Logger
	TaskStartTopic() topic.Topic[bus.TaskStartEvent]
	TaskService() service.TaskService
	GethService() service.GethService
	VandalService() service.VandalService
	ContractService() service.ContractService
}
