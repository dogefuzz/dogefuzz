package service

import (
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/repo"
)

type ContractService interface {
	Create(ctr *dto.NewContractDTO) (*dto.ContractDTO, error)
	Deploy(ctr *common.Contract)
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

func (s *contractService) Create(ctr *dto.NewContractDTO) (*dto.ContractDTO, error) {
	contract := s.contractMapper.ToDomain(ctr)
	err := s.contractRepo.Create(contract)
	if err != nil {
		return nil, err
	}
	contractDTO := s.contractMapper.ToDTO(contract)
	return contractDTO, nil
}
