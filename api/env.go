package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/controller"
	"github.com/dogefuzz/dogefuzz/data"
	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"go.uber.org/zap"
)

type Env interface {
	Logger() *zap.Logger
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
	TasksController() interfaces.TasksController
	TransactionsController() interfaces.TransactionsController
	InstrumentExecutionTopic() interfaces.Topic[bus.InstrumentExecutionEvent]
	TaskFinishTopic() interfaces.Topic[bus.TaskFinishEvent]
	TaskInputRequestTopic() interfaces.Topic[bus.TaskInputRequestEvent]
	TaskStartTopic() interfaces.Topic[bus.TaskStartEvent]
	Deployer() interfaces.Deployer
	Agent() interfaces.Agent
}

type env struct {
	cfg                      *config.Config
	logger                   *zap.Logger
	dbConnection             interfaces.Connection
	eventBus                 interfaces.EventBus
	solidityCompiler         interfaces.SolidityCompiler
	contractMapper           interfaces.ContractMapper
	transactionMapper        interfaces.TransactionMapper
	taskMapper               interfaces.TaskMapper
	functionMapper           interfaces.FunctionMapper
	taskRepo                 interfaces.TaskRepo
	transactionRepo          interfaces.TransactionRepo
	contractRepo             interfaces.ContractRepo
	functionRepo             interfaces.FunctionRepo
	contractService          interfaces.ContractService
	transactionService       interfaces.TransactionService
	taskService              interfaces.TaskService
	functionService          interfaces.FunctionService
	tasksController          interfaces.TasksController
	transactionsController   interfaces.TransactionsController
	instrumentExecutionTopic interfaces.Topic[bus.InstrumentExecutionEvent]
	taskFinishTopic          interfaces.Topic[bus.TaskFinishEvent]
	taskInputRequestTopic    interfaces.Topic[bus.TaskInputRequestEvent]
	taskStartTopic           interfaces.Topic[bus.TaskStartEvent]
	deployer                 interfaces.Deployer
	agent                    interfaces.Agent
}

func NewEnv(cfg *config.Config) *env {
	return &env{cfg: cfg}
}

func (e *env) Destroy() {
	if e.dbConnection != nil {
		e.dbConnection.Clean()
	}
	if e.dbConnection != nil {
		e.logger.Sync()
	}
}

func (e *env) Config() *config.Config {
	return e.cfg
}

func (e *env) Logger() *zap.Logger {
	if e.logger == nil {
		logger, err := initLogger()
		if err != nil {
			log.Panicf("Error while loading zap logger: %s", err)
			return nil
		}

		e.logger = logger
	}
	return e.logger
}

func (e *env) DbConnection() interfaces.Connection {
	if e.dbConnection == nil {
		dbConnection, err := data.NewConnection(e.cfg, e.logger)
		if err != nil {
			e.logger.Error(fmt.Sprintf("Error while initializing database manager: %s", err))
			return nil
		}
		e.dbConnection = dbConnection
	}
	return e.dbConnection
}

func (e *env) EventBus() interfaces.EventBus {
	if e.eventBus == nil {
		e.eventBus = bus.NewMemoryEventBus()
	}
	return e.eventBus
}

func (e *env) SolidityCompiler() interfaces.SolidityCompiler {
	if e.solidityCompiler == nil {
		e.solidityCompiler = solc.NewSolidityCompiler(e.cfg.StorageFolder)
	}
	return e.solidityCompiler
}

func (e *env) ContractMapper() interfaces.ContractMapper {
	if e.contractMapper == nil {
		e.contractMapper = mapper.NewContractMapper()
	}
	return e.contractMapper
}

func (e *env) TransactionMapper() interfaces.TransactionMapper {
	if e.transactionMapper == nil {
		e.transactionMapper = mapper.NewTransactionMapper()
	}
	return e.transactionMapper
}

func (e *env) TaskMapper() interfaces.TaskMapper {
	if e.taskMapper == nil {
		e.taskMapper = mapper.NewTaskMapper()
	}
	return e.taskMapper
}

func (e *env) FunctionMapper() interfaces.FunctionMapper {
	if e.functionMapper == nil {
		e.functionMapper = mapper.NewFunctionMapper()
	}
	return e.functionMapper
}

func (e *env) TaskRepo() interfaces.TaskRepo {
	if e.taskRepo == nil {
		e.taskRepo = repo.NewTaskRepo(e)
	}
	return e.taskRepo
}

func (e *env) TransactionRepo() interfaces.TransactionRepo {
	if e.transactionRepo == nil {
		e.transactionRepo = repo.NewTransactionRepo(e)
	}
	return e.transactionRepo
}

func (e *env) ContractRepo() interfaces.ContractRepo {
	if e.contractRepo == nil {
		e.contractRepo = repo.NewContractRepo(e)
	}
	return e.contractRepo
}

func (e *env) FunctionRepo() interfaces.FunctionRepo {
	if e.functionRepo == nil {
		e.functionRepo = repo.NewFunctionRepo(e)
	}
	return e.functionRepo
}

func (e *env) ContractService() interfaces.ContractService {
	if e.contractService == nil {
		e.contractService = service.NewContractService(e)
	}
	return e.contractService
}

func (e *env) TransactionService() interfaces.TransactionService {
	if e.transactionService == nil {
		e.transactionService = service.NewTransactionService(e)
	}
	return e.transactionService
}

func (e *env) TaskService() interfaces.TaskService {
	if e.taskService == nil {
		e.taskService = service.NewTaskService(e)
	}
	return e.taskService
}

func (e *env) FunctionService() interfaces.FunctionService {
	if e.functionService == nil {
		e.functionService = service.NewFunctionService(e)
	}
	return e.functionService
}

func (e *env) TasksController() interfaces.TasksController {
	if e.tasksController == nil {
		e.tasksController = controller.NewTasksController(e)
	}
	return e.tasksController
}

func (e *env) TransactionsController() interfaces.TransactionsController {
	if e.transactionsController == nil {
		e.transactionsController = controller.NewTransactionsController(e)
	}
	return e.transactionsController
}

func (e *env) InstrumentExecutionTopic() interfaces.Topic[bus.InstrumentExecutionEvent] {
	if e.instrumentExecutionTopic == nil {
		e.instrumentExecutionTopic = topic.NewInstrumentExecutionTopic(e)
	}
	return e.instrumentExecutionTopic
}

func (e *env) TaskFinishTopic() interfaces.Topic[bus.TaskFinishEvent] {
	if e.taskFinishTopic == nil {
		e.taskFinishTopic = topic.NewTaskFinishTopic(e)
	}
	return e.taskFinishTopic
}

func (e *env) TaskInputRequestTopic() interfaces.Topic[bus.TaskInputRequestEvent] {
	if e.taskInputRequestTopic == nil {
		e.taskInputRequestTopic = topic.NewTaskInputRequestTopic(e)
	}
	return e.taskInputRequestTopic
}

func (e *env) TaskStartTopic() interfaces.Topic[bus.TaskStartEvent] {
	if e.taskStartTopic == nil {
		e.taskStartTopic = topic.NewTaskStartTopic(e)
	}
	return e.taskStartTopic
}

func (e *env) Deployer() interfaces.Deployer {
	if e.deployer == nil {
		deployer, err := geth.NewDeployer(e.cfg.GethConfig)
		if err != nil {
			panic(err)
		}
		e.deployer = deployer
	}
	return e.deployer
}

func (e *env) Agent() interfaces.Agent {
	if e.agent == nil {
		agent, err := geth.NewAgent(e.cfg.GethConfig)
		if err != nil {
			panic(err)
		}
		e.agent = agent
	}
	return e.agent
}

func initLogger() (*zap.Logger, error) {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "/tmp/logs"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, err
	}
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return l, nil
}
