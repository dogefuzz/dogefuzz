package mocks

import (
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/stretchr/testify/mock"
)

type ContractMapperMock struct {
	mock.Mock
}

func (m *ContractMapperMock) ToDomain(d *dto.NewContractDTO) *entities.Contract {
	args := m.Called(d)
	return args.Get(0).(*entities.Contract)
}

func (m *ContractMapperMock) ToDTO(c *entities.Contract) *dto.ContractDTO {
	args := m.Called(c)
	return args.Get(0).(*dto.ContractDTO)
}

func (m *ContractMapperMock) ToCommon(c *dto.ContractDTO) *common.Contract {
	args := m.Called(c)
	return args.Get(0).(*common.Contract)
}
