package test

import (
	"encoding/json"
	"log"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/data"
	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
	"go.uber.org/zap"
)

type TestEnv struct {
	cfg               *config.Config
	logger            *zap.Logger
	contractMapper    mapper.ContractMapper
	transactionMapper mapper.TransactionMapper
	taskMapper        mapper.TaskMapper
	functionMapper    mapper.FunctionMapper
	taskRepo          repo.TaskRepo
	contractRepo      repo.ContractRepo
	transactionRepo   repo.TransactionRepo
	functionRepo      repo.FunctionRepo
	deployer          geth.Deployer
	agent             geth.Agent
	connection        data.Connection
}

func NewTestEnv(
	contractMapper mapper.ContractMapper,
	transactionMapper mapper.TransactionMapper,
	taskMapper mapper.TaskMapper,
	functionMapper mapper.FunctionMapper,
	taskRepo repo.TaskRepo,
	contractRepo repo.ContractRepo,
	transactionRepo repo.TransactionRepo,
	functionRepo repo.FunctionRepo,
	deployer geth.Deployer,
	agent geth.Agent,
) *TestEnv {
	return &TestEnv{
		contractMapper: contractMapper,
		contractRepo:   contractRepo,
		deployer:       deployer,
		agent:          agent,
	}
}

func (e *TestEnv) Config() *config.Config {
	return e.cfg
}

func (e *TestEnv) ContractMapper() mapper.ContractMapper {
	return e.contractMapper
}

func (e *TestEnv) TransactionMapper() mapper.TransactionMapper {
	return e.transactionMapper
}

func (e *TestEnv) TaskMapper() mapper.TaskMapper {
	return e.taskMapper
}

func (e *TestEnv) FunctionMapper() mapper.FunctionMapper {
	return e.functionMapper
}

func (e *TestEnv) TaskRepo() repo.TaskRepo {
	return e.taskRepo
}

func (e *TestEnv) TransactionRepo() repo.TransactionRepo {
	return e.transactionRepo
}

func (e *TestEnv) ContractRepo() repo.ContractRepo {
	return e.contractRepo
}

func (e *TestEnv) FunctionRepo() repo.FunctionRepo {
	return e.functionRepo
}

func (e *TestEnv) Agent() geth.Agent {
	return e.agent
}

func (e *TestEnv) DbConnection() data.Connection {
	return e.connection
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

func (e *TestEnv) Deployer() geth.Deployer {
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
