package api

import (
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type Env interface {
	Logger() *zap.Logger
	Client() interfaces.HttpClient
	Config() *config.Config
	DbConnection() interfaces.Connection
	EventBus() interfaces.EventBus
	SolidityCompiler() interfaces.SolidityCompiler
	ContractMapper() interfaces.ContractMapper
	TransactionMapper() interfaces.TransactionMapper
	TaskMapper() interfaces.TaskMapper
	FunctionMapper() interfaces.FunctionMapper
	TaskRepo() interfaces.TaskRepo
	TransactionRepo() interfaces.TransactionRepo
	ContractRepo() interfaces.ContractRepo
	FunctionRepo() interfaces.FunctionRepo
	ContractService() interfaces.ContractService
	TransactionService() interfaces.TransactionService
	TaskService() interfaces.TaskService
	FunctionService() interfaces.FunctionService
	ContractsController() interfaces.ContractsController
	TasksController() interfaces.TasksController
	TransactionsController() interfaces.TransactionsController
	InstrumentExecutionTopic() interfaces.Topic[bus.InstrumentExecutionEvent]
	TaskFinishTopic() interfaces.Topic[bus.TaskFinishEvent]
	TaskInputRequestTopic() interfaces.Topic[bus.TaskInputRequestEvent]
	TaskStartTopic() interfaces.Topic[bus.TaskStartEvent]
	Deployer() interfaces.Deployer
	Agent() interfaces.Agent
}
