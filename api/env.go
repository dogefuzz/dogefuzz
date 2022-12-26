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
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"go.uber.org/zap"
)

type Env interface {
	Logger() *zap.Logger
	DbConnection() data.Connection
	EventBus() bus.EventBus
	SolidityCompiler() solc.SolidityCompiler
	ContractMapper() mapper.ContractMapper
	TransactionMapper() mapper.TransactionMapper
	TaskMapper() mapper.TaskMapper
	FunctionMapper() mapper.FunctionMapper
	TaskRepo() repo.TaskRepo
	TransactionRepo() repo.TransactionRepo
	ContractRepo() repo.ContractRepo
	FunctionRepo() repo.FunctionRepo
	ContractService() service.ContractService
	TransactionService() service.TransactionService
	TaskService() service.TaskService
	FunctionService() service.FunctionService
	TasksController() controller.TasksController
	TransactionsController() controller.TransactionsController
	InstrumentExecutionTopic() topic.Topic[bus.InstrumentExecutionEvent]
	TaskFinishTopic() topic.Topic[bus.TaskFinishEvent]
	TaskInputRequestTopic() topic.Topic[bus.TaskInputRequestEvent]
	TaskStartTopic() topic.Topic[bus.TaskStartEvent]
	Deployer() geth.Deployer
}

type env struct {
	cfg                      *config.Config
	logger                   *zap.Logger
	dbConnection             data.Connection
	eventBus                 bus.EventBus
	solidityCompiler         solc.SolidityCompiler
	contractMapper           mapper.ContractMapper
	transactionMapper        mapper.TransactionMapper
	taskMapper               mapper.TaskMapper
	functionMapper           mapper.FunctionMapper
	taskRepo                 repo.TaskRepo
	transactionRepo          repo.TransactionRepo
	contractRepo             repo.ContractRepo
	functionRepo             repo.FunctionRepo
	contractService          service.ContractService
	transactionService       service.TransactionService
	taskService              service.TaskService
	functionService          service.FunctionService
	tasksController          controller.TasksController
	transactionsController   controller.TransactionsController
	instrumentExecutionTopic topic.Topic[bus.InstrumentExecutionEvent]
	taskFinishTopic          topic.Topic[bus.TaskFinishEvent]
	taskInputRequestTopic    topic.Topic[bus.TaskInputRequestEvent]
	taskStartTopic           topic.Topic[bus.TaskStartEvent]
	deployer                 geth.Deployer
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

func (e *env) DbConnection() data.Connection {
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

func (e *env) EventBus() bus.EventBus {
	if e.eventBus == nil {
		e.eventBus = bus.NewMemoryEventBus()
	}
	return e.eventBus
}

func (e *env) SolidityCompiler() solc.SolidityCompiler {
	if e.solidityCompiler == nil {
		e.solidityCompiler = solc.NewSolidityCompiler(e.cfg.StorageFolder)
	}
	return e.solidityCompiler
}

func (e *env) ContractMapper() mapper.ContractMapper {
	if e.contractMapper == nil {
		e.contractMapper = mapper.NewContractMapper()
	}
	return e.contractMapper
}

func (e *env) TransactionMapper() mapper.TransactionMapper {
	if e.transactionMapper == nil {
		e.transactionMapper = mapper.NewTransactionMapper()
	}
	return e.transactionMapper
}

func (e *env) TaskMapper() mapper.TaskMapper {
	if e.taskMapper == nil {
		e.taskMapper = mapper.NewTaskMapper()
	}
	return e.taskMapper
}

func (e *env) FunctionMapper() mapper.FunctionMapper {
	if e.functionMapper == nil {
		e.functionMapper = mapper.NewFunctionMapper()
	}
	return e.functionMapper
}

func (e *env) TaskRepo() repo.TaskRepo {
	if e.taskRepo == nil {
		e.taskRepo = repo.NewTaskRepo(e)
	}
	return e.taskRepo
}

func (e *env) TransactionRepo() repo.TransactionRepo {
	if e.transactionRepo == nil {
		e.transactionRepo = repo.NewTransactionRepo(e)
	}
	return e.transactionRepo
}

func (e *env) ContractRepo() repo.ContractRepo {
	if e.contractRepo == nil {
		e.contractRepo = repo.NewContractRepo(e)
	}
	return e.contractRepo
}

func (e *env) FunctionRepo() repo.FunctionRepo {
	if e.functionRepo == nil {
		e.functionRepo = repo.NewFunctionRepo(e)
	}
	return e.functionRepo
}

func (e *env) ContractService() service.ContractService {
	if e.contractService == nil {
		e.contractService = service.NewContractService(e)
	}
	return e.contractService
}

func (e *env) TransactionService() service.TransactionService {
	if e.transactionService == nil {
		e.transactionService = service.NewTransactionService(e)
	}
	return e.transactionService
}

func (e *env) TaskService() service.TaskService {
	if e.taskService == nil {
		e.taskService = service.NewTaskService(e)
	}
	return e.taskService
}

func (e *env) FunctionService() service.FunctionService {
	if e.functionService == nil {
		e.functionService = service.NewFunctionService(e)
	}
	return e.functionService
}

func (e *env) TasksController() controller.TasksController {
	if e.tasksController == nil {
		e.tasksController = controller.NewTasksController(e)
	}
	return e.tasksController
}

func (e *env) TransactionsController() controller.TransactionsController {
	if e.transactionsController == nil {
		e.transactionsController = controller.NewTransactionsController(e)
	}
	return e.transactionsController
}

func (e *env) InstrumentExecutionTopic() topic.Topic[bus.InstrumentExecutionEvent] {
	if e.instrumentExecutionTopic == nil {
		e.instrumentExecutionTopic = topic.NewInstrumentExecutionTopic(e)
	}
	return e.instrumentExecutionTopic
}

func (e *env) TaskFinishTopic() topic.Topic[bus.TaskFinishEvent] {
	if e.taskFinishTopic == nil {
		e.taskFinishTopic = topic.NewTaskFinishTopic(e)
	}
	return e.taskFinishTopic
}

func (e *env) TaskInputRequestTopic() topic.Topic[bus.TaskInputRequestEvent] {
	if e.taskInputRequestTopic == nil {
		e.taskInputRequestTopic = topic.NewTaskInputRequestTopic(e)
	}
	return e.taskInputRequestTopic
}

func (e *env) TaskStartTopic() topic.Topic[bus.TaskStartEvent] {
	if e.taskStartTopic == nil {
		e.taskStartTopic = topic.NewTaskStartTopic(e)
	}
	return e.taskStartTopic
}

func (e *env) Deployer() geth.Deployer {
	if e.deployer == nil {
		deployer, err := geth.NewDeployer(e.cfg.GethConfig)
		if err != nil {
			panic(err)
		}
		e.deployer = deployer
	}
	return e.deployer
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
