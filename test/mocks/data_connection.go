package mocks

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ConnectionMock struct {
	mock.Mock
}

func (m *ConnectionMock) Clean() error {
	args := m.Called()
	return args.Error(0)
}

func (m *ConnectionMock) GetDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *ConnectionMock) Migrate() error {
	args := m.Called()
	return args.Error(0)
}
