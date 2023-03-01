package mocks

import (
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ContractRepoMock struct {
	mock.Mock
}

func (m *ContractRepoMock) Create(tx *gorm.DB, contract *entities.Contract) error {
	args := m.Called(tx, contract)
	return args.Error(0)
}

func (m *ContractRepoMock) Update(tx *gorm.DB, contract *entities.Contract) error {
	args := m.Called(tx, contract)
	return args.Error(0)
}

func (m *ContractRepoMock) FindAll(tx *gorm.DB) ([]entities.Contract, error) {
	args := m.Called(tx)
	return args.Get(0).([]entities.Contract), args.Error(1)
}

func (m *ContractRepoMock) Find(tx *gorm.DB, id string) (*entities.Contract, error) {
	args := m.Called(tx, id)
	return args.Get(0).(*entities.Contract), args.Error(1)
}

func (m *ContractRepoMock) FindByTaskId(tx *gorm.DB, taskId string) (*entities.Contract, error) {
	args := m.Called(tx, taskId)
	return args.Get(0).(*entities.Contract), args.Error(1)
}
