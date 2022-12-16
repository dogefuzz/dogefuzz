package mapper

import (
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type ContractMapper interface {
	ToDomain(d *dto.NewContractDTO) *domain.Contract
	ToDTO(c *domain.Contract) *dto.ContractDTO
	ToCommon(c *dto.ContractDTO) *common.Contract
}

type contractMapper struct{}

func NewContractMapper() *contractMapper {
	return &contractMapper{}
}

func (m *contractMapper) ToDomain(d *dto.NewContractDTO) *domain.Contract {
	return &domain.Contract{
		Name:   d.Name,
		Source: d.Source,
	}
}

func (m *contractMapper) ToDTO(c *domain.Contract) *dto.ContractDTO {
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
	}
}
