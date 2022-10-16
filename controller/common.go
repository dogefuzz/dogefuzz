package controller

import (
	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/db"
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
}
