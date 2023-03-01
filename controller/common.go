package controller

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type Env interface {
	Logger() *zap.Logger

	ContractService() interfaces.ContractService
	TransactionService() interfaces.TransactionService
	TaskService() interfaces.TaskService
	FunctionService() interfaces.FunctionService
	SolidityService() interfaces.SolidityService

	SolidityCompiler() interfaces.SolidityCompiler

	InstrumentExecutionTopic() interfaces.Topic[bus.InstrumentExecutionEvent]
	TaskStartTopic() interfaces.Topic[bus.TaskStartEvent]
}
