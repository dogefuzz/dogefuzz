package service

import (
	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/repo"
)

type Env interface {
	ContractMapper() mapper.ContractMapper
	ContractRepo() repo.ContractRepo
}
