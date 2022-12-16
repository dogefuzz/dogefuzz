package mapper

import (
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/google/uuid"
)

type FunctionMapper interface {
	ToDomainForCreation(c *dto.NewFunctionDTO) *domain.Function
	ToDomain(c *dto.FunctionDTO) *domain.Function
	ToDTO(c *domain.Function) *dto.FunctionDTO
}

type functionMapper struct{}

func NewFunctionMapper() *functionMapper {
	return &functionMapper{}
}

func (m *functionMapper) ToDomainForCreation(c *dto.NewFunctionDTO) *domain.Function {
	return &domain.Function{
		Id:           uuid.NewString(),
		Name:         c.Name,
		NumberOfArgs: c.NumberOfArgs,
	}
}

func (m *functionMapper) ToDomain(c *dto.FunctionDTO) *domain.Function {
	return &domain.Function{
		Id:           c.Id,
		Name:         c.Name,
		NumberOfArgs: c.NumberOfArgs,
	}
}

func (m *functionMapper) ToDTO(c *domain.Function) *dto.FunctionDTO {
	return &dto.FunctionDTO{
		Name:         c.Name,
		NumberOfArgs: c.NumberOfArgs,
	}
}
