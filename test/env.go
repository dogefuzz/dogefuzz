package test

import (
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type TestEnv struct {
	cfg                    *config.Config
	logger                 *zap.Logger
	contractMapper         interfaces.ContractMapper
	transactionMapper      interfaces.TransactionMapper
	taskMapper             interfaces.TaskMapper
	functionMapper         interfaces.FunctionMapper
	taskRepo               interfaces.TaskRepo
	contractRepo           interfaces.ContractRepo
	transactionRepo        interfaces.TransactionRepo
	functionRepo           interfaces.FunctionRepo
	deployer               interfaces.Deployer
	agent                  interfaces.Agent
	connection             interfaces.Connection
	eventBus               interfaces.EventBus
	taskService            interfaces.TaskService
	taskInputRequestTopic  interfaces.Topic[bus.TaskInputRequestEvent]
	taskFinishTopic        interfaces.Topic[bus.TaskFinishEvent]
	tasksCheckerJob        interfaces.CronJob
	transactionsCheckerJob interfaces.CronJob
}

func NewTestEnv(
	contractMapper interfaces.ContractMapper,
	transactionMapper interfaces.TransactionMapper,
	taskMapper interfaces.TaskMapper,
	functionMapper interfaces.FunctionMapper,
	taskRepo interfaces.TaskRepo,
	contractRepo interfaces.ContractRepo,
	transactionRepo interfaces.TransactionRepo,
	functionRepo interfaces.FunctionRepo,
	deployer interfaces.Deployer,
	agent interfaces.Agent,
	connection interfaces.Connection,
	eventBus interfaces.EventBus,
	taskService interfaces.TaskService,
	taskFinishTopic interfaces.Topic[bus.TaskFinishEvent],
	taskInputRequestTopic interfaces.Topic[bus.TaskInputRequestEvent],
) *TestEnv {
	return &TestEnv{
		contractMapper:        contractMapper,
		contractRepo:          contractRepo,
		deployer:              deployer,
		agent:                 agent,
		connection:            connection,
		eventBus:              eventBus,
		taskService:           taskService,
		taskFinishTopic:       taskFinishTopic,
		taskInputRequestTopic: taskInputRequestTopic,
	}
}

func (e *TestEnv) Config() *config.Config {
	return e.cfg
}

func (e *TestEnv) ContractMapper() interfaces.ContractMapper {
	return e.contractMapper
}

func (e *TestEnv) TransactionMapper() interfaces.TransactionMapper {
	return e.transactionMapper
}

func (e *TestEnv) TaskMapper() interfaces.TaskMapper {
	return e.taskMapper
}

func (e *TestEnv) FunctionMapper() interfaces.FunctionMapper {
	return e.functionMapper
}

func (e *TestEnv) TaskRepo() interfaces.TaskRepo {
	return e.taskRepo
}

func (e *TestEnv) TransactionRepo() interfaces.TransactionRepo {
	return e.transactionRepo
}

func (e *TestEnv) ContractRepo() interfaces.ContractRepo {
	return e.contractRepo
}

func (e *TestEnv) FunctionRepo() interfaces.FunctionRepo {
	return e.functionRepo
}

func (e *TestEnv) TaskService() interfaces.TaskService {
	return e.taskService
}

func (e *TestEnv) TaskInputRequestTopic() interfaces.Topic[bus.TaskInputRequestEvent] {
	return e.taskInputRequestTopic
}

func (e *TestEnv) TaskFinishTopic() interfaces.Topic[bus.TaskFinishEvent] {
	return e.taskFinishTopic
}

func (e *TestEnv) TasksCheckerJob() interfaces.CronJob {
	return e.tasksCheckerJob
}

func (e *TestEnv) TransactionsCheckerJob() interfaces.CronJob {
	return e.transactionsCheckerJob
}

func (e *TestEnv) Agent() interfaces.Agent {
	return e.agent
}

func (e *TestEnv) Deployer() interfaces.Deployer {
	return e.deployer
}

func (e *TestEnv) DbConnection() interfaces.Connection {
	return e.connection
}

func (e *TestEnv) EventBus() interfaces.EventBus {
	return e.eventBus
}

func (e *TestEnv) Logger() *zap.Logger {
	if e.logger == nil {
		e.logger = zap.NewNop()
	}
	return e.logger
}
