package service

import (
	"context"

	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
	"go.uber.org/zap"
)

type Env interface {
	ContractMapper() mapper.ContractMapper
	TransactionMapper() mapper.TransactionMapper
	TaskMapper() mapper.TaskMapper
	FunctionMapper() mapper.FunctionMapper
	TaskRepo() repo.TaskRepo
	ContractRepo() repo.ContractRepo
	TransactionRepo() repo.TransactionRepo
	FunctionRepo() repo.FunctionRepo
	Logger() *zap.Logger
	Deployer() geth.Deployer
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
