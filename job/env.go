package job

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

	ContractMapper() interfaces.ContractMapper
	FunctionMapper() interfaces.FunctionMapper
	TaskMapper() interfaces.TaskMapper
	TransactionMapper() interfaces.TransactionMapper

	ContractRepo() interfaces.ContractRepo
	FunctionRepo() interfaces.FunctionRepo
	TaskRepo() interfaces.TaskRepo
	TransactionRepo() interfaces.TransactionRepo

	TaskService() interfaces.TaskService

	TaskInputRequestTopic() interfaces.Topic[bus.TaskInputRequestEvent]
	TaskFinishTopic() interfaces.Topic[bus.TaskFinishEvent]

	TasksCheckerJob() interfaces.CronJob
	TransactionsCheckerJob() interfaces.CronJob
	Deployer() interfaces.Deployer
	Agent() interfaces.Agent
}
