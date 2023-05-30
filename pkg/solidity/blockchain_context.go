package solidity

type BlockchainContext struct {
	AvailableAddresses []string
}

func NewBlockchainContext(availableAddresses []string) *BlockchainContext {
	return &BlockchainContext{
		AvailableAddresses: availableAddresses,
	}
}
