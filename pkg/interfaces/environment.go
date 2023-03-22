package interfaces

import "context"

type BlockchainEnvironment interface {
	Setup(ctx context.Context) error
}

type ContractPool interface {
	Setup(ctx context.Context) error
}
