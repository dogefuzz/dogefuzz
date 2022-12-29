package service

import (
	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
)

type ContractService interface {
	Get(contractId string) (*dto.ContractDTO, error)
	FindByTaskId(contractId string) (*dto.ContractDTO, error)
	Create(ctr *dto.NewContractDTO) (*dto.ContractDTO, error)
	Update(ctr *dto.ContractDTO) error
}

type contractService struct {
	contractMapper mapper.ContractMapper
	contractRepo   repo.ContractRepo
}

func NewContractService(e Env) *contractService {
	return &contractService{
		contractMapper: e.ContractMapper(),
		contractRepo:   e.ContractRepo(),
	}
}

func (s *contractService) Get(contractId string) (*dto.ContractDTO, error) {
	// TODO: get contract by ID
	return nil, nil
}

func (s *contractService) FindByTaskId(contractId string) (*dto.ContractDTO, error) {
	// TODO: get contract by task id
	return nil, nil
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

func (s *contractService) Update(ctr *dto.ContractDTO) error {
	// TODO: Add update logic
	return nil
}
