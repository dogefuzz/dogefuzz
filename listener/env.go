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
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"go.uber.org/zap"
)

type Env interface {
	Config() *config.Config
	Logger() *zap.Logger
	TaskStartTopic() interfaces.Topic[bus.TaskStartEvent]
	TaskFinishTopic() interfaces.Topic[bus.TaskFinishEvent]
	TaskInputRequestTopic() interfaces.Topic[bus.TaskInputRequestEvent]
	InstrumentExecutionTopic() interfaces.Topic[bus.InstrumentExecutionEvent]
	TaskService() interfaces.TaskService
	GethService() interfaces.GethService
	VandalService() interfaces.VandalService
	ContractService() interfaces.ContractService
	FunctionService() interfaces.FunctionService
	ContractMapper() interfaces.ContractMapper
	TransactionService() interfaces.TransactionService
	ReporterService() interfaces.ReporterService
	ContractDeployerListener() interfaces.Listener
	ExecutionAnalyticsListener() interfaces.Listener
	FuzzerListener() interfaces.Listener
	ReporterListener() interfaces.Listener
	Deployer() interfaces.Deployer
	Agent() interfaces.Agent
	FuzzerLeader() interfaces.FuzzerLeader
	BlackboxFuzzer() interfaces.Fuzzer
	GreyboxFuzzer() interfaces.Fuzzer
	DirectedGreyboxFuzzer() interfaces.Fuzzer
	PowerSchedule() interfaces.PowerSchedule
}

type env struct {
	cfg                        *config.Config
	logger                     *zap.Logger
	dbConnection               interfaces.Connection
	eventBus                   bus.EventBus
	contractMapper             interfaces.ContractMapper
	functionMapper             interfaces.FunctionMapper
	taskMapper                 interfaces.TaskMapper
	transactionMapper          interfaces.TransactionMapper
	contractRepo               interfaces.ContractRepo
	functionRepo               interfaces.FunctionRepo
	taskRepo                   interfaces.TaskRepo
	transactionRepo            interfaces.TransactionRepo
	taskService                interfaces.TaskService
	contractService            interfaces.ContractService
	functionService            interfaces.FunctionService
	transactionService         interfaces.TransactionService
	gethService                interfaces.GethService
	vandalService              interfaces.VandalService
	reporterService            interfaces.ReporterService
	taskInputRequestTopic      interfaces.Topic[bus.TaskInputRequestEvent]
	taskStartTopic             interfaces.Topic[bus.TaskStartEvent]
	taskFinishTopic            interfaces.Topic[bus.TaskFinishEvent]
	instrumentExecutionTopic   interfaces.Topic[bus.InstrumentExecutionEvent]
	contractDeployerListener   interfaces.Listener
	executionAnalyticsListener interfaces.Listener
	fuzzerListener             interfaces.Listener
	reporterListener           interfaces.Listener
	deployer                   interfaces.Deployer
	agent                      interfaces.Agent
	fuzzerLeader               interfaces.FuzzerLeader
	blackboxFuzzer             interfaces.Fuzzer
	greyboxFuzzer              interfaces.Fuzzer
	directedGreyboxFuzzer      interfaces.Fuzzer
	powerSchedule              interfaces.PowerSchedule
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

func (e *env) EventBus() bus.EventBus {
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

func (e *env) ContractService() interfaces.ContractService {
	if e.contractService == nil {
		e.contractService = service.NewContractService(e)
	}
	return e.contractService
}

func (e *env) TaskService() interfaces.TaskService {
	if e.taskService == nil {
		e.taskService = service.NewTaskService(e)
	}
	return e.taskService
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

func (e *env) FunctionService() interfaces.FunctionService {
	if e.functionService == nil {
		e.functionService = service.NewFunctionService(e)
	}
	return e.functionService
}

func (e *env) TransactionService() interfaces.TransactionService {
	if e.transactionService == nil {
		e.transactionService = service.NewTransactionService(e)
	}
	return e.transactionService
}

func (e *env) ReporterService() interfaces.ReporterService {
	if e.reporterService == nil {
		e.reporterService = service.NewReporterService(e)
	}
	return e.reporterService
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

func (e *env) TaskFinishTopic() interfaces.Topic[bus.TaskFinishEvent] {
	if e.taskFinishTopic == nil {
		e.taskFinishTopic = topic.NewTaskFinishTopic(e)
	}
	return e.taskFinishTopic
}

func (e *env) InstrumentExecutionTopic() interfaces.Topic[bus.InstrumentExecutionEvent] {
	if e.instrumentExecutionTopic == nil {
		e.instrumentExecutionTopic = topic.NewInstrumentExecutionTopic(e)
	}
	return e.instrumentExecutionTopic
}

func (e *env) ContractDeployerListener() interfaces.Listener {
	if e.contractDeployerListener == nil {
		e.contractDeployerListener = NewContractDeployerListener(e)
	}
	return e.contractDeployerListener
}

func (e *env) ExecutionAnalyticsListener() interfaces.Listener {
	if e.executionAnalyticsListener == nil {
		e.executionAnalyticsListener = NewExecutionAnalyticsListener(e)
	}
	return e.executionAnalyticsListener
}

func (e *env) FuzzerListener() interfaces.Listener {
	if e.fuzzerListener == nil {
		e.fuzzerListener = NewFuzzerListener(e)
	}
	return e.fuzzerListener
}

func (e *env) ReporterListener() interfaces.Listener {
	if e.reporterListener == nil {
		e.reporterListener = NewReporterListener(e)
	}
	return e.reporterListener
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

func (e *env) FuzzerLeader() interfaces.FuzzerLeader {
	if e.fuzzerLeader == nil {
		e.fuzzerLeader = fuzz.NewFuzzerLeader(e)
	}
	return e.fuzzerLeader
}

func (e *env) BlackboxFuzzer() interfaces.Fuzzer {
	if e.blackboxFuzzer == nil {
		e.blackboxFuzzer = fuzz.NewBlackboxFuzzer()
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
