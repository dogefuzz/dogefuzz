package mapper

import (
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type ContractMapper interface {
	ToDomain(d *dto.NewContractDTO) *entities.Contract
	ToDTO(c *entities.Contract) *dto.ContractDTO
	ToCommon(c *dto.ContractDTO) *common.Contract
}

type contractMapper struct{}

func NewContractMapper() *contractMapper {
	return &contractMapper{}
}

func (m *contractMapper) ToDomain(d *dto.NewContractDTO) *entities.Contract {
	return &entities.Contract{
		Name:   d.Name,
		Source: d.Source,
	}
}

func (m *contractMapper) ToDTO(c *entities.Contract) *dto.ContractDTO {
	return &dto.ContractDTO{
		Id:      c.Id,
		Address: c.Address,
		Source:  c.Source,
		Name:    c.Name,
	}
}

func (m *contractMapper) ToCommon(c *dto.ContractDTO) *common.Contract {
	return &common.Contract{
		Name:          c.Name,
		AbiDefinition: c.AbiDefinition,
		CompiledCode:  c.CompiledCode,
		Address:       c.Address,
	}
}
