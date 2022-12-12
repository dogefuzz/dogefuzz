package service

import (
	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/repo"
)

type Env interface {
	ContractMapper() mapper.ContractMapper
	TransactionMapper() mapper.TransactionMapper
	TaskMapper() mapper.TaskMapper
	OracleMapper() mapper.OracleMapper
	TaskOracleRepo() repo.TaskOracleRepo
	TaskRepo() repo.TaskRepo
	ContractRepo() repo.ContractRepo
	TransactionRepo() repo.TransactionRepo
	OracleRepo() repo.OracleRepo
	Deployer() geth.Deployer
}
