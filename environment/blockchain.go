package environment

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

type blockchainEnvironment struct {
	contractPool interfaces.ContractPool
}

func NewBlockchainEnvironment(e env) *blockchainEnvironment {
	return &blockchainEnvironment{
		contractPool: e.ContractPool(),
	}
}

func (e *blockchainEnvironment) Setup(ctx context.Context) error {
	return e.contractPool.Setup(ctx)
}
