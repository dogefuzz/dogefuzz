package service

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
)

var ErrFunctionNotFound = errors.New("function not found")

type FunctionService interface {
	Get(functionId string) (*dto.FunctionDTO, error)
	Create(task *dto.NewFunctionDTO) (*dto.FunctionDTO, error)
	FindByContractId(contractId string) []*dto.FunctionDTO
}

type functionService struct {
	functionRepo   repo.FunctionRepo
	functionMapper mapper.FunctionMapper
}

func NewFunctionService(e Env) *functionService {
	return &functionService{
		functionRepo:   e.FunctionRepo(),
		functionMapper: e.FunctionMapper(),
	}
}

func (s *functionService) Get(functionId string) (*dto.FunctionDTO, error) {
	transaction, err := s.functionRepo.GetById(functionId)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrFunctionNotFound
		}
		return nil, err
	}
	return s.functionMapper.ToDTO(transaction), nil
}

func (s *functionService) Create(function *dto.NewFunctionDTO) (*dto.FunctionDTO, error) {
	entity := s.functionMapper.ToDomainForCreation(function)
	err := s.functionRepo.Create(entity)
	if err != nil {
		return nil, err
	}
	return s.functionMapper.ToDTO(entity), nil
}

func (s *functionService) FindByContractId(contractId string) []*dto.FunctionDTO {
	// TODO: call repository to find functions of some specific contract
	return nil
}
