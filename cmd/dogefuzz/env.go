package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/controller"
	"github.com/dogefuzz/dogefuzz/data"
	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/environment"
	"github.com/dogefuzz/dogefuzz/fuzz"
	"github.com/dogefuzz/dogefuzz/job"
	"github.com/dogefuzz/dogefuzz/listener"
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
	Client() interfaces.HttpClient
	Config() *config.Config
	DbConnection() interfaces.Connection
	EventBus() interfaces.EventBus
	SolidityCompiler() interfaces.SolidityCompiler
	Deployer() interfaces.Deployer
	Agent() interfaces.Agent
	ContractPool() interfaces.ContractPool

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
	GethService() interfaces.GethService
	VandalService() interfaces.VandalService
	ReporterService() interfaces.ReporterService
	SolidityService() interfaces.SolidityService

	ContractsController() interfaces.ContractsController
	TasksController() interfaces.TasksController
	TransactionsController() interfaces.TransactionsController
	PingController() interfaces.PingController

	InstrumentExecutionTopic() interfaces.Topic[bus.InstrumentExecutionEvent]
	TaskFinishTopic() interfaces.Topic[bus.TaskFinishEvent]
	TaskInputRequestTopic() interfaces.Topic[bus.TaskInputRequestEvent]
	TaskStartTopic() interfaces.Topic[bus.TaskStartEvent]

	ContractDeployerListener() interfaces.Listener
	ExecutionAnalyticsListener() interfaces.Listener
	FuzzerListener() interfaces.Listener
	ReporterListener() interfaces.Listener

	TasksCheckerJob() interfaces.CronJob
	TransactionsCheckerJob() interfaces.CronJob
	TransactionsTimeoutCheckerJob() interfaces.CronJob

	FuzzerLeader() interfaces.FuzzerLeader
	BlackboxFuzzer() interfaces.Fuzzer
	GreyboxFuzzer() interfaces.Fuzzer
	DirectedGreyboxFuzzer() interfaces.Fuzzer
	PowerSchedule() interfaces.PowerSchedule
}

type env struct {
	cfg              *config.Config
	logger           *zap.Logger
	client           interfaces.HttpClient
	dbConnection     interfaces.Connection
	eventBus         interfaces.EventBus
	solidityCompiler interfaces.SolidityCompiler
	deployer         interfaces.Deployer
	agent            interfaces.Agent
	contractPool     interfaces.ContractPool

	contractMapper    interfaces.ContractMapper
	transactionMapper interfaces.TransactionMapper
	taskMapper        interfaces.TaskMapper
	functionMapper    interfaces.FunctionMapper

	taskRepo        interfaces.TaskRepo
	transactionRepo interfaces.TransactionRepo
	contractRepo    interfaces.ContractRepo
	functionRepo    interfaces.FunctionRepo

	contractService    interfaces.ContractService
	transactionService interfaces.TransactionService
	taskService        interfaces.TaskService
	functionService    interfaces.FunctionService
	gethService        interfaces.GethService
	vandalService      interfaces.VandalService
	reporterService    interfaces.ReporterService
	solidityService    interfaces.SolidityService

	contractsController    interfaces.ContractsController
	tasksController        interfaces.TasksController
	transactionsController interfaces.TransactionsController
	pingController         interfaces.PingController

	instrumentExecutionTopic interfaces.Topic[bus.InstrumentExecutionEvent]
	taskFinishTopic          interfaces.Topic[bus.TaskFinishEvent]
	taskInputRequestTopic    interfaces.Topic[bus.TaskInputRequestEvent]
	taskStartTopic           interfaces.Topic[bus.TaskStartEvent]

	contractDeployerListener   interfaces.Listener
	executionAnalyticsListener interfaces.Listener
	fuzzerListener             interfaces.Listener
	reporterListener           interfaces.Listener

	tasksCheckerJob               interfaces.CronJob
	transactionsCheckerJob        interfaces.CronJob
	transactionsTimeoutCheckerJob interfaces.CronJob

	fuzzerLeader          interfaces.FuzzerLeader
	blackboxFuzzer        interfaces.Fuzzer
	greyboxFuzzer         interfaces.Fuzzer
	directedGreyboxFuzzer interfaces.Fuzzer
	powerSchedule         interfaces.PowerSchedule
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

func (e *env) Client() interfaces.HttpClient {
	if e.client == nil {
		e.client = new(http.Client)
	}
	return e.client
}

func (e *env) DbConnection() interfaces.Connection {
	if e.dbConnection == nil {
		dbConnection, err := data.NewConnection(e.cfg, e.logger)
		if err != nil {
			e.logger.Error(fmt.Sprintf("Error while initializing database manager: %s", err))
			panic(err)
		}
		err = dbConnection.Migrate()
		if err != nil {
			e.logger.Error(fmt.Sprintf("Error while migrating database: %s", err))
			panic(err)
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
		e.solidityCompiler = solc.NewSolidityCompiler(e.Config().StorageFolder)
	}
	return e.solidityCompiler
}

func (e *env) Deployer() interfaces.Deployer {
	if e.deployer == nil {
		deployer, err := geth.NewDeployer(e.Logger(), e.Config().GethConfig)
		if err != nil {
			panic(err)
		}
		e.deployer = deployer
	}
	return e.deployer
}

func (e *env) Agent() interfaces.Agent {
	if e.agent == nil {
		agent, err := geth.NewAgent(e.Config().GethConfig)
		if err != nil {
			panic(err)
		}
		e.agent = agent
	}
	return e.agent
}

func (e *env) ContractPool() interfaces.ContractPool {
	if e.contractPool == nil {
		e.contractPool = environment.NewContractPool(e)
	}
	return e.contractPool
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

func (e *env) GethService() interfaces.GethService {
	if e.gethService == nil {
		e.gethService = service.NewGethService(e)
	}
	return e.gethService
}

func (e *env) VandalService() interfaces.VandalService {
	if e.vandalService == nil {
		e.vandalService = service.NewVandalService(e)
	}
	return e.vandalService
}

func (e *env) ReporterService() interfaces.ReporterService {
	if e.reporterService == nil {
		e.reporterService = service.NewReporterService(e)
	}
	return e.reporterService
}

func (e *env) SolidityService() interfaces.SolidityService {
	if e.solidityService == nil {
		e.solidityService = service.NewSolidityService(e)
	}
	return e.solidityService
}

func (e *env) ContractsController() interfaces.ContractsController {
	if e.contractsController == nil {
		e.contractsController = controller.NewContractsController(e)
	}
	return e.contractsController
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

func (e *env) PingController() interfaces.PingController {
	if e.pingController == nil {
		e.pingController = controller.NewPingController()
	}
	return e.pingController
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

func (e *env) ContractDeployerListener() interfaces.Listener {
	if e.contractDeployerListener == nil {
		e.contractDeployerListener = listener.NewContractDeployerListener(e)
	}
	return e.contractDeployerListener
}

func (e *env) ExecutionAnalyticsListener() interfaces.Listener {
	if e.executionAnalyticsListener == nil {
		e.executionAnalyticsListener = listener.NewExecutionAnalyticsListener(e)
	}
	return e.executionAnalyticsListener
}

func (e *env) FuzzerListener() interfaces.Listener {
	if e.fuzzerListener == nil {
		e.fuzzerListener = listener.NewFuzzerListener(e)
	}
	return e.fuzzerListener
}

func (e *env) ReporterListener() interfaces.Listener {
	if e.reporterListener == nil {
		e.reporterListener = listener.NewReporterListener(e)
	}
	return e.reporterListener
}

func (e *env) TasksCheckerJob() interfaces.CronJob {
	if e.tasksCheckerJob == nil {
		e.tasksCheckerJob = job.NewTasksCheckerJob(e)
	}
	return e.tasksCheckerJob
}

func (e *env) TransactionsCheckerJob() interfaces.CronJob {
	if e.transactionsCheckerJob == nil {
		e.transactionsCheckerJob = job.NewTransactionsCheckerJob(e)
	}
	return e.transactionsCheckerJob
}

func (e *env) TransactionsTimeoutCheckerJob() interfaces.CronJob {
	if e.transactionsTimeoutCheckerJob == nil {
		e.transactionsTimeoutCheckerJob = job.NewTransactionsTimeoutCheckerJob(e)
	}
	return e.transactionsTimeoutCheckerJob
}

func (e *env) FuzzerLeader() interfaces.FuzzerLeader {
	if e.fuzzerLeader == nil {
		e.fuzzerLeader = fuzz.NewFuzzerLeader(e)
	}
	return e.fuzzerLeader
}

func (e *env) BlackboxFuzzer() interfaces.Fuzzer {
	if e.blackboxFuzzer == nil {
		e.blackboxFuzzer = fuzz.NewBlackboxFuzzer(e)
	}
	return e.blackboxFuzzer
}

func (e *env) GreyboxFuzzer() interfaces.Fuzzer {
	if e.greyboxFuzzer == nil {
		e.greyboxFuzzer = fuzz.NewGreyboxFuzzer(e)
	}
	return e.greyboxFuzzer
}

func (e *env) DirectedGreyboxFuzzer() interfaces.Fuzzer {
	if e.directedGreyboxFuzzer == nil {
		e.directedGreyboxFuzzer = fuzz.NewDirectedGreyboxFuzzer(e)
	}
	return e.directedGreyboxFuzzer
}

func (e *env) PowerSchedule() interfaces.PowerSchedule {
	if e.powerSchedule == nil {
		e.powerSchedule = fuzz.NewPowerSchedule(e)
	}
	return e.powerSchedule
}

func initLogger() (*zap.Logger, error) {
	// rawJSON := []byte(`{
	// 	"level": "debug",
	// 	"encoding": "console",
	// 	"outputPaths": ["stdout", "/tmp/logs"],
	// 	"errorOutputPaths": ["stderr"],
	// 	"encoderConfig": {
	// 		"messageKey": "message",
	// 		"levelKey": "level",
	// 		"levelEncoder": "lowercase"
	// 	}
	// }`)

	// var cfg zap.Config
	// if err := json.Unmarshal(rawJSON, &cfg); err != nil {
	// 	return nil, err
	// }
	// l, err := cfg.
	// 	WithTimestampFormat("2006-01-02 15:04:05.000").
	// 	Build()
	l, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		return nil, err
	}
	return l, nil
}
