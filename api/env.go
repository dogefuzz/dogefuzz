package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/controller"
	"github.com/dogefuzz/dogefuzz/db"
	"github.com/dogefuzz/dogefuzz/repo"
	"go.uber.org/zap"
)

type Env interface {
	Logger() *zap.Logger
	DbConnection() db.Connection
	EventBus() bus.EventBus
	OracleRepo() repo.OracleRepo
	TaskOracleRepo() repo.TaskOracleRepo
	TaskRepo() repo.TaskRepo
	TransactionRepo() repo.TransactionRepo
	ContractRepo() repo.ContractRepo
	TaskContractRepo() repo.TaskContractRepo
	TasksController() controller.TasksController
	WeaknessesController() controller.WeaknessesController
	ExecutionsController() controller.ExecutionsController
	TransactionsController() controller.TransactionsController
}

type env struct {
	cfg                    *config.Config
	logger                 *zap.Logger
	dbConnection           db.Connection
	eventBus               bus.EventBus
	oracleRepo             repo.OracleRepo
	taskOracleRepo         repo.TaskOracleRepo
	taskRepo               repo.TaskRepo
	transactionRepo        repo.TransactionRepo
	contractRepo           repo.ContractRepo
	taskContractRepo       repo.TaskContractRepo
	tasksController        controller.TasksController
	weaknessesController   controller.WeaknessesController
	executionsController   controller.ExecutionsController
	transactionsController controller.TransactionsController
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

func (e *env) DbConnection() db.Connection {
	if e.dbConnection == nil {
		dbConnection, err := db.NewConnection(e.cfg, e.logger)
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

func (e *env) OracleRepo() repo.OracleRepo {
	if e.oracleRepo == nil {
		e.oracleRepo = repo.NewOracleRepo(e)
	}
	return e.oracleRepo
}

func (e *env) TaskOracleRepo() repo.TaskOracleRepo {
	if e.taskOracleRepo == nil {
		e.taskOracleRepo = repo.NewTaskOracleRepo(e)
	}
	return e.taskOracleRepo
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

func (e *env) TaskContractRepo() repo.TaskContractRepo {
	if e.taskContractRepo == nil {
		e.taskContractRepo = repo.NewTaskContractRepo(e)
	}
	return e.taskContractRepo
}

func (e *env) ExecutionsController() controller.ExecutionsController {
	if e.executionsController == nil {
		e.executionsController = controller.NewExecutionsController(e)
	}
	return e.executionsController
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

func (e *env) WeaknessesController() controller.WeaknessesController {
	if e.weaknessesController == nil {
		e.weaknessesController = controller.NewWeaknessesController(e)
	}
	return e.weaknessesController
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
