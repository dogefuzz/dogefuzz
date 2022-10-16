package mapper

import (
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
)

type ContractMapper interface {
	ToDomain(d *dto.NewContractDTO) *domain.Contract
	ToDTO(c *domain.Contract) *dto.ContractDTO
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
