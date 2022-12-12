package mapper

import (
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
)

type OracleMapper interface {
	ToDomain(c *dto.OracleDTO) *domain.Oracle
	ToDTO(c *domain.Oracle) *dto.OracleDTO
}

type oracleMapper struct{}

func NewOracleMapper() *oracleMapper {
	return &oracleMapper{}
}

func (m *oracleMapper) ToDomain(c *dto.OracleDTO) *domain.Oracle {
	return &domain.Oracle{
		Id:   c.Id,
		Name: c.Name,
	}
}

func (m *oracleMapper) ToDTO(c *domain.Oracle) *dto.OracleDTO {
	return &dto.OracleDTO{
		Id:   c.Id,
		Name: c.Name,
	}
}
