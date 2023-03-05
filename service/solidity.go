package service

import (
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/pkg/solidity"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type solidityService struct {
	cfg          config.GethConfig
	contractRepo interfaces.ContractRepo
	connection   interfaces.Connection
}

func NewSolidityService(e Env) *solidityService {
	return &solidityService{
		cfg:          e.Config().GethConfig,
		contractRepo: e.ContractRepo(),
		connection:   e.DbConnection(),
	}
}

func (s *solidityService) GetTypeHandlerWithContext(typ abi.Type) (interfaces.TypeHandler, error) {
	blockchainContext := &solidity.BlockchainContext{}
	addresses, err := s.getAvailableAddresses()
	if err != nil {
		return nil, err
	}
	blockchainContext.AvailableAddresses = addresses
	return solidity.GetTypeHandler(typ, blockchainContext)
}

func (s *solidityService) getAvailableAddresses() ([]string, error) {
	availableAddresses := make([]string, 0)
	if s.cfg.AgentAddress != "" {
		availableAddresses = append(availableAddresses, s.cfg.AgentAddress)
	}

	if s.cfg.DeployerAddress != "" {
		availableAddresses = append(availableAddresses, s.cfg.DeployerAddress)
	}

	contracts, err := s.contractRepo.FindAll(s.connection.GetDB())
	if err != nil {
		return nil, err
	}

	for _, contract := range contracts {
		if contract.Address != "" {
			availableAddresses = append(availableAddresses, contract.Address)
		}
	}
	return availableAddresses, nil
}
