package listener

import (
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/fuzz"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"go.uber.org/zap"
)

type env interface {
	Config() *config.Config
	Logger() *zap.Logger
	FuzzerLeader() fuzz.FuzzerLeader
	TaskStartTopic() topic.Topic[bus.TaskStartEvent]
	TaskInputRequestTopic() topic.Topic[bus.TaskInputRequestEvent]
	TaskService() service.TaskService
	GethService() service.GethService
	VandalService() service.VandalService
	ContractService() service.ContractService
	FunctionService() service.FunctionService
	ContractMapper() mapper.ContractMapper
	TransactionService() service.TransactionService
}
