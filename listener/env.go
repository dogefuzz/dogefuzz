package listener

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/data"
	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/fuzz"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"go.uber.org/zap"
)

type Env interface {
	Config() *config.Config
	Logger() *zap.Logger
	TaskStartTopic() topic.Topic[bus.TaskStartEvent]
	TaskFinishTopic() topic.Topic[bus.TaskFinishEvent]
	TaskInputRequestTopic() topic.Topic[bus.TaskInputRequestEvent]
	InstrumentExecutionTopic() topic.Topic[bus.InstrumentExecutionEvent]
	TaskService() service.TaskService
	GethService() service.GethService
	VandalService() service.VandalService
	ContractService() service.ContractService
	FunctionService() service.FunctionService
	ContractMapper() mapper.ContractMapper
	TransactionService() service.TransactionService
	ReporterService() service.ReporterService
	ContractDeployerListener() Listener
	ExecutionAnalyticsListener() Listener
	FuzzerListener() Listener
	ReporterListener() Listener
	Deployer() geth.Deployer
	Agent() geth.Agent
	FuzzerLeader() fuzz.FuzzerLeader
	BlackboxFuzzer() fuzz.Fuzzer
	GreyboxFuzzer() fuzz.Fuzzer
	DirectedGreyboxFuzzer() fuzz.Fuzzer
	PowerSchedule() fuzz.PowerSchedule
}

type env struct {
	cfg                        *config.Config
	logger                     *zap.Logger
	dbConnection               data.Connection
	eventBus                   bus.EventBus
	contractMapper             mapper.ContractMapper
	functionMapper             mapper.FunctionMapper
	taskMapper                 mapper.TaskMapper
	transactionMapper          mapper.TransactionMapper
	contractRepo               repo.ContractRepo
	functionRepo               repo.FunctionRepo
	taskRepo                   repo.TaskRepo
	transactionRepo            repo.TransactionRepo
	taskService                service.TaskService
	contractService            service.ContractService
	functionService            service.FunctionService
	transactionService         service.TransactionService
	gethService                service.GethService
	vandalService              service.VandalService
	reporterService            service.ReporterService
	taskInputRequestTopic      topic.Topic[bus.TaskInputRequestEvent]
	taskStartTopic             topic.Topic[bus.TaskStartEvent]
	taskFinishTopic            topic.Topic[bus.TaskFinishEvent]
	instrumentExecutionTopic   topic.Topic[bus.InstrumentExecutionEvent]
	contractDeployerListener   Listener
	executionAnalyticsListener Listener
	fuzzerListener             Listener
	reporterListener           Listener
	deployer                   geth.Deployer
	agent                      geth.Agent
	fuzzerLeader               fuzz.FuzzerLeader
	blackboxFuzzer             fuzz.Fuzzer
	greyboxFuzzer              fuzz.Fuzzer
	directedGreyboxFuzzer      fuzz.Fuzzer
	powerSchedule              fuzz.PowerSchedule
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

func (e *env) ContractMapper() mapper.ContractMapper {
	if e.contractMapper == nil {
		e.contractMapper = mapper.NewContractMapper()
	}
	return e.contractMapper
}

func (e *env) FunctionMapper() mapper.FunctionMapper {
	if e.functionMapper == nil {
		e.functionMapper = mapper.NewFunctionMapper()
	}
	return e.functionMapper
}

func (e *env) TaskMapper() mapper.TaskMapper {
	if e.taskMapper == nil {
		e.taskMapper = mapper.NewTaskMapper()
	}
	return e.taskMapper
}

func (e *env) TransactionMapper() mapper.TransactionMapper {
	if e.transactionMapper == nil {
		e.transactionMapper = mapper.NewTransactionMapper()
	}
	return e.transactionMapper
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

func (e *env) ContractService() service.ContractService {
	if e.contractService == nil {
		e.contractService = service.NewContractService(e)
	}
	return e.contractService
}

func (e *env) TaskService() service.TaskService {
	if e.taskService == nil {
		e.taskService = service.NewTaskService(e)
	}
	return e.taskService
}

func (e *env) GethService() service.GethService {
	if e.gethService == nil {
		e.gethService = service.NewGethService(e)
	}
	return e.gethService
}

func (e *env) VandalService() service.VandalService {
	if e.vandalService == nil {
		e.vandalService = service.NewVandalService(e)
	}
	return e.vandalService
}

func (e *env) FunctionService() service.FunctionService {
	if e.functionService == nil {
		e.functionService = service.NewFunctionService(e)
	}
	return e.functionService
}

func (e *env) TransactionService() service.TransactionService {
	if e.transactionService == nil {
		e.transactionService = service.NewTransactionService(e)
	}
	return e.transactionService
}

func (e *env) ReporterService() service.ReporterService {
	if e.reporterService == nil {
		e.reporterService = service.NewReporterService(e)
	}
	return e.reporterService
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

func (e *env) TaskFinishTopic() topic.Topic[bus.TaskFinishEvent] {
	if e.taskFinishTopic == nil {
		e.taskFinishTopic = topic.NewTaskFinishTopic(e)
	}
	return e.taskFinishTopic
}

func (e *env) InstrumentExecutionTopic() topic.Topic[bus.InstrumentExecutionEvent] {
	if e.instrumentExecutionTopic == nil {
		e.instrumentExecutionTopic = topic.NewInstrumentExecutionTopic(e)
	}
	return e.instrumentExecutionTopic
}

func (e *env) ContractDeployerListener() Listener {
	if e.contractDeployerListener == nil {
		e.contractDeployerListener = NewContractDeployerListener(e)
	}
	return e.contractDeployerListener
}

func (e *env) ExecutionAnalyticsListener() Listener {
	if e.executionAnalyticsListener == nil {
		e.executionAnalyticsListener = NewExecutionAnalyticsListener(e)
	}
	return e.executionAnalyticsListener
}

func (e *env) FuzzerListener() Listener {
	if e.fuzzerListener == nil {
		e.fuzzerListener = NewFuzzerListener(e)
	}
	return e.fuzzerListener
}

func (e *env) ReporterListener() Listener {
	if e.reporterListener == nil {
		e.reporterListener = NewReporterListener(e)
	}
	return e.reporterListener
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

func (e *env) Agent() geth.Agent {
	if e.agent == nil {
		agent, err := geth.NewAgent(e.cfg.GethConfig)
		if err != nil {
			panic(err)
		}
		e.agent = agent
	}
	return e.agent
}

func (e *env) FuzzerLeader() fuzz.FuzzerLeader {
	if e.fuzzerLeader == nil {
		e.fuzzerLeader = fuzz.NewFuzzerLeader(e)
	}
	return e.fuzzerLeader
}

func (e *env) BlackboxFuzzer() fuzz.Fuzzer {
	if e.blackboxFuzzer == nil {
		e.blackboxFuzzer = fuzz.NewBlackboxFuzzer()
	}
	return e.blackboxFuzzer
}

func (e *env) GreyboxFuzzer() fuzz.Fuzzer {
	if e.greyboxFuzzer == nil {
		e.greyboxFuzzer = fuzz.NewGreyboxFuzzer(e)
	}
	return e.greyboxFuzzer
}

func (e *env) DirectedGreyboxFuzzer() fuzz.Fuzzer {
	if e.directedGreyboxFuzzer == nil {
		e.directedGreyboxFuzzer = fuzz.NewDirectedGreyboxFuzzer(e)
	}
	return e.directedGreyboxFuzzer
}

func (e *env) PowerSchedule() fuzz.PowerSchedule {
	if e.powerSchedule == nil {
		e.powerSchedule = fuzz.NewPowerSchedule(e)
	}
	return e.powerSchedule
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
