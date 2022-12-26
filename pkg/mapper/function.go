package mapper

import (
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/google/uuid"
)

type FunctionMapper interface {
	ToDomainForCreation(c *dto.NewFunctionDTO) *entities.Function
	ToDomain(c *dto.FunctionDTO) *entities.Function
	ToDTO(c *entities.Function) *dto.FunctionDTO
}

type functionMapper struct{}

func NewFunctionMapper() *functionMapper {
	return &functionMapper{}
}

func (m *functionMapper) ToDomainForCreation(c *dto.NewFunctionDTO) *entities.Function {
	return &entities.Function{
		Id:           uuid.NewString(),
		Name:         c.Name,
		NumberOfArgs: c.NumberOfArgs,
	}
}

func (m *functionMapper) ToDomain(c *dto.FunctionDTO) *entities.Function {
	return &entities.Function{
		Id:           c.Id,
		Name:         c.Name,
		NumberOfArgs: c.NumberOfArgs,
	}
}

func (m *functionMapper) ToDTO(c *entities.Function) *dto.FunctionDTO {
	return &dto.FunctionDTO{
		Name:         c.Name,
		NumberOfArgs: c.NumberOfArgs,
	}
}
