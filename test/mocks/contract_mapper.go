package mocks

import (
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/stretchr/testify/mock"
)

type ContractMapperMock struct {
	mock.Mock
}

func (m *ContractMapperMock) ToDomain(d *dto.NewContractDTO) *domain.Contract {
	args := m.Called(d)
	return args.Get(0).(*domain.Contract)
}

func (m *ContractMapperMock) ToDTO(c *domain.Contract) *dto.ContractDTO {
	args := m.Called(c)
	return args.Get(0).(*dto.ContractDTO)
}
