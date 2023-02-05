package mapper

import (
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type functionMapper struct{}

func NewFunctionMapper() *functionMapper {
	return &functionMapper{}
}

func (m *functionMapper) MapNewDTOToEntity(c *dto.NewFunctionDTO) *entities.Function {
	return &entities.Function{
		Name:                    c.Name,
		NumberOfArgs:            c.NumberOfArgs,
		IsChangingContractState: c.IsChangingContractState,
		IsConstructor:           c.IsConstructor,
		ContractId:              c.ContractId,
	}
}

func (m *functionMapper) MapDTOToEntity(c *dto.FunctionDTO) *entities.Function {
	return &entities.Function{
		Id:                      c.Id,
		Name:                    c.Name,
		NumberOfArgs:            c.NumberOfArgs,
		IsChangingContractState: c.IsChangingContractState,
		IsConstructor:           c.IsConstructor,
		ContractId:              c.ContractId,
	}
}

func (m *functionMapper) MapEntityToDTO(c *entities.Function) *dto.FunctionDTO {
	return &dto.FunctionDTO{
		Name:                    c.Name,
		NumberOfArgs:            c.NumberOfArgs,
		IsChangingContractState: c.IsChangingContractState,
		IsConstructor:           c.IsConstructor,
		ContractId:              c.ContractId,
	}
}
