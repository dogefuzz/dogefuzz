package environment

import (
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type env interface {
	Logger() *zap.Logger
	ContractPool() interfaces.ContractPool
	GethService() interfaces.GethService
	ContractService() interfaces.ContractService
	TransactionService() interfaces.TransactionService
	SolidityCompiler() interfaces.SolidityCompiler
}
