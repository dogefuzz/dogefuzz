package controller

import (
	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/dogefuzz/dogefuzz/service"
	"go.uber.org/zap"
)

type Env interface {
	Logger() *zap.Logger
	ContractService() service.ContractService
	TransactionService() service.TransactionService
	TaskService() service.TaskService
	FunctionService() service.FunctionService
	SolidityCompiler() solc.SolidityCompiler
	InstrumentExecutionTopic() topic.Topic[bus.InstrumentExecutionEvent]
	TaskStartTopic() topic.Topic[bus.TaskStartEvent]
}
