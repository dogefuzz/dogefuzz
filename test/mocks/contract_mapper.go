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

func (m *ContractMapperMock) MapNewDTOToEntity(d *dto.NewContractDTO) *entities.Contract {
	args := m.Called(d)
	return args.Get(0).(*entities.Contract)
}

func (m *ContractMapperMock) MapNewWithIdDTOToEntity(d *dto.NewContractWithIdDTO) *entities.Contract {
	args := m.Called(d)
	return args.Get(0).(*entities.Contract)
}

func (m *ContractMapperMock) MapDTOToEntity(d *dto.ContractDTO) *entities.Contract {
	args := m.Called(d)
	return args.Get(0).(*entities.Contract)
}

func (m *ContractMapperMock) MapEntityToDTO(c *entities.Contract) *dto.ContractDTO {
	args := m.Called(c)
	return args.Get(0).(*dto.ContractDTO)
}

func (m *ContractMapperMock) MapDTOToCommon(d *dto.ContractDTO) *common.Contract {
	args := m.Called(d)
	return args.Get(0).(*common.Contract)
}
