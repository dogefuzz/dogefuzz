package job

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/data"
	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"go.uber.org/zap"
)

type Env interface {
	Config() *config.Config
	Logger() *zap.Logger
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

type env struct {
	cfg                    *config.Config
	logger                 *zap.Logger
	dbConnection           interfaces.Connection
	eventBus               interfaces.EventBus
	contractMapper         interfaces.ContractMapper
	functionMapper         interfaces.FunctionMapper
	taskMapper             interfaces.TaskMapper
	transactionMapper      interfaces.TransactionMapper
	contractRepo           interfaces.ContractRepo
	functionRepo           interfaces.FunctionRepo
	taskRepo               interfaces.TaskRepo
	transactionRepo        interfaces.TransactionRepo
	taskService            interfaces.TaskService
	taskInputRequestTopic  interfaces.Topic[bus.TaskInputRequestEvent]
	taskFinishTopic        interfaces.Topic[bus.TaskFinishEvent]
	tasksCheckerJob        interfaces.CronJob
	transactionsCheckerJob interfaces.CronJob
	deployer               interfaces.Deployer
	agent                  interfaces.Agent
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

func (e *env) ContractMapper() interfaces.ContractMapper {
	if e.contractMapper == nil {
		e.contractMapper = mapper.NewContractMapper()
	}
	return e.contractMapper
}

func (e *env) FunctionMapper() interfaces.FunctionMapper {
	if e.functionMapper == nil {
		e.functionMapper = mapper.NewFunctionMapper()
	}
	return e.functionMapper
}

func (e *env) TaskMapper() interfaces.TaskMapper {
	if e.taskMapper == nil {
		e.taskMapper = mapper.NewTaskMapper()
	}
	return e.taskMapper
}

func (e *env) TransactionMapper() interfaces.TransactionMapper {
	if e.transactionMapper == nil {
		e.transactionMapper = mapper.NewTransactionMapper()
	}
	return e.transactionMapper
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

func (e *env) TaskService() interfaces.TaskService {
	if e.taskService == nil {
		e.taskService = service.NewTaskService(e)
	}
	return e.taskService
}

func (e *env) TaskInputRequestTopic() interfaces.Topic[bus.TaskInputRequestEvent] {
	if e.taskInputRequestTopic == nil {
		e.taskInputRequestTopic = topic.NewTaskInputRequestTopic(e)
	}
	return e.taskInputRequestTopic
}

func (e *env) TaskFinishTopic() interfaces.Topic[bus.TaskFinishEvent] {
	if e.taskFinishTopic == nil {
		e.taskFinishTopic = topic.NewTaskFinishTopic(e)
	}
	return e.taskFinishTopic
}

func (e *env) TasksCheckerJob() interfaces.CronJob {
	if e.tasksCheckerJob == nil {
		e.tasksCheckerJob = NewTasksCheckerJob(e)
	}
	return e.tasksCheckerJob
}

func (e *env) TransactionsCheckerJob() interfaces.CronJob {
	if e.transactionsCheckerJob == nil {
		e.transactionsCheckerJob = NewTransactionsCheckerJob(e)
	}
	return e.transactionsCheckerJob
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
