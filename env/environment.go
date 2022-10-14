package env

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gongbell/contractfuzzer/bus"
	"github.com/gongbell/contractfuzzer/db"
	"github.com/gongbell/contractfuzzer/db/repository"
	"go.uber.org/zap"
)

type Environment interface {
	Logger() *zap.Logger
	DbManager() db.Manager
	EventBus() bus.EventBus
	OracleRepository() repository.OracleRepository
	TaskOracleRepository() repository.TaskOracleRepository
	TaskRepository() repository.TaskRepository
	TransactionRepository() repository.TransactionRepository
	ContractRepository() repository.ContractRepository
	TaskContractRepository() repository.TaskContractRepository
}

type DefaultEnvironment struct {
	logger                 *zap.Logger
	dbManager              db.Manager
	eventBus               bus.EventBus
	oracleRepository       repository.OracleRepository
	taskOracleRepository   repository.TaskOracleRepository
	taskRepository         repository.TaskRepository
	transactionRepository  repository.TransactionRepository
	contractRepository     repository.ContractRepository
	taskContractRepository repository.TaskContractRepository
}

func (e DefaultEnvironment) Init() (DefaultEnvironment, error) {

	// Init logger
	logger, err := initLogger()
	if err != nil {
		log.Panicf("Error while loading zap logger: %s", err)
		return DefaultEnvironment{}, err
	}
	e.logger = logger
	e.logger.Info("Log framework initialized with success")

	// Init DB manager
	dbManager, err := new(db.SQLiteManager).Init(e.logger)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error while initializing database manager: %s", err))
		return DefaultEnvironment{}, err
	}
	e.dbManager = dbManager
	e.logger.Info("Database manager initialized with success")

	// Init event bus
	eventBus, err := new(bus.MemoryEventBus).Init()
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error while initializing event bus: %s", err))
		return DefaultEnvironment{}, err
	}
	e.eventBus = eventBus
	e.logger.Info("Event bus initialized with success")

	// Init oracle repository
	oracleRepository, err := new(repository.OracleSQLiteRepository).Init(e.dbManager)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error while initializing oracle repository: %s", err))
		return DefaultEnvironment{}, err
	}
	e.oracleRepository = oracleRepository

	// Init task_oracle repository
	taskOracleRepository, err := new(repository.TaskOracleSQLiteRepository).Init(e.dbManager)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error while initializing task_oracle repository: %s", err))
		return DefaultEnvironment{}, err
	}
	e.taskOracleRepository = taskOracleRepository

	// Init task repository
	taskRepository, err := new(repository.TaskSQLiteRepository).Init(e.dbManager)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error while initializing task repository: %s", err))
		return DefaultEnvironment{}, err
	}
	e.taskRepository = taskRepository

	// Init transaction repository
	transactionRepository, err := new(repository.TransactionSQLiteRepository).Init(e.dbManager)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error while initializing transaction repository: %s", err))
		return DefaultEnvironment{}, err
	}
	e.transactionRepository = transactionRepository

	// Init contract repository
	contractRepository, err := new(repository.ContractSQLiteRepository).Init(e.dbManager)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error while initializing contract repository: %s", err))
		return DefaultEnvironment{}, err
	}
	e.contractRepository = contractRepository

	// Init task_contract repository
	taskContractRepository, err := new(repository.TaskContractSQLiteRepository).Init(e.dbManager)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error while initializing task_contract repository: %s", err))
		return DefaultEnvironment{}, err
	}
	e.taskContractRepository = taskContractRepository
	e.logger.Info("Domain repositories initialized with success")

	e.logger.Info("Environment initialized with success")
	return e, nil
}

func (e DefaultEnvironment) Destroy() {
	e.dbManager.Clean()
	e.logger.Sync()
}

func (e DefaultEnvironment) Logger() *zap.Logger {
	return e.logger
}

func (e DefaultEnvironment) DbManager() db.Manager {
	return e.dbManager
}

func (e DefaultEnvironment) EventBus() bus.EventBus {
	return e.eventBus
}

func (e DefaultEnvironment) OracleRepository() repository.OracleRepository {
	return e.oracleRepository
}

func (e DefaultEnvironment) TaskOracleRepository() repository.TaskOracleRepository {
	return e.taskOracleRepository
}

func (e DefaultEnvironment) TaskRepository() repository.TaskRepository {
	return e.taskRepository
}

func (e DefaultEnvironment) TransactionRepository() repository.TransactionRepository {
	return e.transactionRepository
}

func (e DefaultEnvironment) ContractRepository() repository.ContractRepository {
	return e.contractRepository
}

func (e DefaultEnvironment) TaskContractRepository() repository.TaskContractRepository {
	return e.taskContractRepository
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
