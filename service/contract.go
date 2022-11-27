package service

import (
	"context"

	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/repo"
)

type ContractService interface {
	Create(ctr *dto.NewContractDTO) (*dto.ContractDTO, error)
	Deploy(ctx context.Context, contract *common.Contract, args ...string) (string, error)
}

type contractService struct {
	contractMapper mapper.ContractMapper
	contractRepo   repo.ContractRepo
	deployer       geth.Deployer
}

func NewContractService(e Env) *contractService {
	return &contractService{
		contractMapper: e.ContractMapper(),
		contractRepo:   e.ContractRepo(),
		deployer:       e.Deployer(),
	}
}

func (s *contractService) Create(ctr *dto.NewContractDTO) (*dto.ContractDTO, error) {
	contract := s.contractMapper.ToDomain(ctr)
	err := s.contractRepo.Create(contract)
	if err != nil {
		return nil, err
	}
	contractDTO := s.contractMapper.ToDTO(contract)
	return contractDTO, nil
}

func (s *contractService) Deploy(ctx context.Context, contract *common.Contract, args ...string) (string, error) {
	address, err := s.deployer.Deploy(ctx, contract, common.ConvertStringArrayToInterfaceArray(args)...)
	if err != nil {
		return "", err
	}
	return address, nil
}
