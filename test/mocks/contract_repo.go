package mocks

import (
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/stretchr/testify/mock"
)

type ContractRepoMock struct {
	mock.Mock
}

func (m *ContractRepoMock) Create(contract *entities.Contract) error {
	args := m.Called(contract)
	return args.Error(0)
}

func (m *ContractRepoMock) Find(id string) (*entities.Contract, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Contract), args.Error(1)
}

func (m *ContractRepoMock) FindByName(name string) (*entities.Contract, error) {
	args := m.Called(name)
	return args.Get(0).(*entities.Contract), args.Error(1)
}

func (m *ContractRepoMock) FindByAddress(address string) (*entities.Contract, error) {
	args := m.Called(address)
	return args.Get(0).(*entities.Contract), args.Error(1)
}

func (m *ContractRepoMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
