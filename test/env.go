package test

import (
	"encoding/json"
	"log"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type TestEnv struct {
	cfg               *config.Config
	logger            *zap.Logger
	contractMapper    interfaces.ContractMapper
	transactionMapper interfaces.TransactionMapper
	taskMapper        interfaces.TaskMapper
	functionMapper    interfaces.FunctionMapper
	taskRepo          interfaces.TaskRepo
	contractRepo      interfaces.ContractRepo
	transactionRepo   interfaces.TransactionRepo
	functionRepo      interfaces.FunctionRepo
	deployer          interfaces.Deployer
	agent             interfaces.Agent
	connection        interfaces.Connection
	eventBus          bus.EventBus
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
	eventBus bus.EventBus,
) *TestEnv {
	return &TestEnv{
		contractMapper: contractMapper,
		contractRepo:   contractRepo,
		deployer:       deployer,
		agent:          agent,
		connection:     connection,
		eventBus:       eventBus,
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

func (e *TestEnv) Agent() interfaces.Agent {
	return e.agent
}

func (e *TestEnv) Deployer() interfaces.Deployer {
	return e.deployer
}

func (e *TestEnv) DbConnection() interfaces.Connection {
	return e.connection
}

func (e *TestEnv) EventBus() bus.EventBus {
	return e.eventBus
}

func (e *TestEnv) Logger() *zap.Logger {
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
