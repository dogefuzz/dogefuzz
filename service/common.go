package service

import (
	"context"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type Env interface {
	ContractMapper() interfaces.ContractMapper
	TransactionMapper() interfaces.TransactionMapper
	TaskMapper() interfaces.TaskMapper
	FunctionMapper() interfaces.FunctionMapper
	TaskRepo() interfaces.TaskRepo
	ContractRepo() interfaces.ContractRepo
	TransactionRepo() interfaces.TransactionRepo
	FunctionRepo() interfaces.FunctionRepo
	Logger() *zap.Logger
	Deployer() interfaces.Deployer
	Agent() interfaces.Agent
	Config() *config.Config
	DbConnection() interfaces.Connection
	Client() interfaces.HttpClient
}

func packArgsToVariadicFuncParameters(ctx context.Context, contract *common.Contract, args []interface{}) []interface{} {
	parameters := make([]interface{}, len(args)+2)
	parameters[0] = ctx
	parameters[1] = contract
	for idx := 0; idx < len(args); idx++ {
		parameters[idx+2] = args[idx]
	}
	return parameters
}
