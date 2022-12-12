package controller

import (
	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/db"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/dogefuzz/dogefuzz/repo"
	"github.com/dogefuzz/dogefuzz/service"
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
	ContractService() service.ContractService
	OracleService() service.OracleService
	TransactionService() service.TransactionService
	TaskService() service.TaskService
	SolidityCompiler() solc.SolidityCompiler
	InstrumentExecutionTopic() topic.Topic[bus.InstrumentExecutionEvent]
	TaskStartTopic() topic.Topic[bus.TaskStartEvent]
}
