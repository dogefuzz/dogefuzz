package service

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

var ErrContractNotFound = errors.New("the contract was not found")

type contractService struct {
	contractMapper interfaces.ContractMapper
	contractRepo   interfaces.ContractRepo
	connection     interfaces.Connection
}

func NewContractService(e Env) *contractService {
	return &contractService{
		contractMapper: e.ContractMapper(),
		contractRepo:   e.ContractRepo(),
		connection:     e.DbConnection(),
	}
}

func (s *contractService) Find(contractId string) (*dto.ContractDTO, error) {
	contract, err := s.contractRepo.Find(s.connection.GetDB(), contractId)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, nil
		}
		return nil, err
	}
	return s.contractMapper.MapEntityToDTO(contract), nil
}

func (s *contractService) Get(contractId string) (*dto.ContractDTO, error) {
	contract, err := s.contractRepo.Find(s.connection.GetDB(), contractId)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrContractNotFound
		}
		return nil, err
	}
	return s.contractMapper.MapEntityToDTO(contract), nil
}

func (s *contractService) FindByTaskId(taskId string) (*dto.ContractDTO, error) {
	contract, err := s.contractRepo.FindByTaskId(s.connection.GetDB(), taskId)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrContractNotFound
		}
		return nil, err
	}
	return s.contractMapper.MapEntityToDTO(contract), nil
}

func (s *contractService) Create(ctr *dto.NewContractDTO) (*dto.ContractDTO, error) {
	contract := s.contractMapper.MapNewDTOToEntity(ctr)
	err := s.contractRepo.Create(s.connection.GetDB(), contract)
	if err != nil {
		return nil, err
	}
	contractDTO := s.contractMapper.MapEntityToDTO(contract)
	return contractDTO, nil
}

func (s *contractService) CreateWithId(ctr *dto.NewContractWithIdDTO) (*dto.ContractDTO, error) {
	contract := s.contractMapper.MapNewWithIdDTOToEntity(ctr)
	err := s.contractRepo.Create(s.connection.GetDB(), contract)
	if err != nil {
		return nil, err
	}
	contractDTO := s.contractMapper.MapEntityToDTO(contract)
	return contractDTO, nil
}

func (s *contractService) Update(ctr *dto.ContractDTO) error {
	contract := s.contractMapper.MapDTOToEntity(ctr)
	err := s.contractRepo.Update(s.connection.GetDB(), contract)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return ErrContractNotFound
		}
		return err
	}
	return nil
}
