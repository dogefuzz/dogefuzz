package service

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

var ErrFunctionNotFound = errors.New("function not found")

type functionService struct {
	connection     interfaces.Connection
	functionRepo   interfaces.FunctionRepo
	functionMapper interfaces.FunctionMapper
}

func NewFunctionService(e Env) *functionService {
	return &functionService{
		functionRepo:   e.FunctionRepo(),
		functionMapper: e.FunctionMapper(),
		connection:     e.DbConnection(),
	}
}

func (s *functionService) Get(functionId string) (*dto.FunctionDTO, error) {
	function, err := s.functionRepo.Get(s.connection.GetDB(), functionId)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrFunctionNotFound
		}
		return nil, err
	}
	return s.functionMapper.MapEntityToDTO(function), nil
}

func (s *functionService) Create(function *dto.NewFunctionDTO) (*dto.FunctionDTO, error) {
	entity := s.functionMapper.MapNewDTOToEntity(function)
	err := s.functionRepo.Create(s.connection.GetDB(), entity)
	if err != nil {
		return nil, err
	}
	return s.functionMapper.MapEntityToDTO(entity), nil
}

func (s *functionService) FindByContractId(contractId string) ([]*dto.FunctionDTO, error) {
	functions, err := s.functionRepo.FindByContractId(s.connection.GetDB(), contractId)
	if err != nil {
		return nil, err
	}
	functionsDTO := make([]*dto.FunctionDTO, len(functions))
	for idx, function := range functions {
		functionsDTO[idx] = s.functionMapper.MapEntityToDTO(&function)
	}
	return functionsDTO, nil
}

func (s *functionService) FindConstructorByContractId(contractId string) (*dto.FunctionDTO, error) {
	function, err := s.functionRepo.FindConstructorByContractId(s.connection.GetDB(), contractId)
	if err != nil {
		return nil, err
	}
	return s.functionMapper.MapEntityToDTO(function), nil
}
