package mocks

import (
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/stretchr/testify/mock"
)

type ContractRepoMock struct {
	mock.Mock
}

func (m *ContractRepoMock) Create(contract *domain.Contract) error {
	args := m.Called(contract)
	return args.Error(0)
}

func (m *ContractRepoMock) Find(id string) (*domain.Contract, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Contract), args.Error(1)
}

func (m *ContractRepoMock) FindByName(name string) (*domain.Contract, error) {
	args := m.Called(name)
	return args.Get(0).(*domain.Contract), args.Error(1)
}

func (m *ContractRepoMock) FindByAddress(address string) (*domain.Contract, error) {
	args := m.Called(address)
	return args.Get(0).(*domain.Contract), args.Error(1)
}

func (m *ContractRepoMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
