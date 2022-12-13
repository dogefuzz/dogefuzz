package service

import (
	"context"

	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/repo"
	"github.com/dogefuzz/dogefuzz/test/mocks"
)

var contractMapperMock = new(mocks.ContractMapperMock)
var contractRepoMock = new(mocks.ContractRepoMock)
var deployerMock = new(mocks.DeployerMock)
var e = &env{
	contractMapper: contractMapperMock,
	contractRepo:   contractRepoMock,
	deployer:       deployerMock,
}

type env struct {
	contractMapper    mapper.ContractMapper
	transactionMapper mapper.TransactionMapper
	taskMapper        mapper.TaskMapper
	taskRepo          repo.TaskRepo
	contractRepo      repo.ContractRepo
	transactionRepo   repo.TransactionRepo
	deployer          geth.Deployer
}

func (e *env) ContractMapper() mapper.ContractMapper {
	return e.contractMapper
}

func (e *env) TransactionMapper() mapper.TransactionMapper {
	return e.transactionMapper
}

func (e *env) TaskMapper() mapper.TaskMapper {
	return e.taskMapper
}

func (e *env) TaskRepo() repo.TaskRepo {
	return e.taskRepo
}

func (e *env) TransactionRepo() repo.TransactionRepo {
	return e.transactionRepo
}

func (e *env) ContractRepo() repo.ContractRepo {
	return e.contractRepo
}

func (e *env) Deployer() geth.Deployer {
	return e.deployer
}

func packArgsToVariadicFuncParameters(ctx context.Context, contract *common.Contract, args []string) []interface{} {
	parameters := make([]interface{}, len(args)+2)
	parameters[0] = ctx
	parameters[1] = contract
	for idx := 0; idx < len(args); idx++ {
		parameters[idx+2] = args[idx]
	}
	return parameters
}
