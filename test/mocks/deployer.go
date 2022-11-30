package mocks

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/stretchr/testify/mock"
)

type DeployerMock struct {
	mock.Mock
}

func (m *DeployerMock) Deploy(ctx context.Context, contract *common.Contract, args ...interface{}) (string, error) {
	arguments := make([]interface{}, 0)
	arguments = append(arguments, ctx)
	arguments = append(arguments, contract)
	if len(args) > 0 {
		arguments = append(arguments, args...)
	}
	mockArgs := m.Called(arguments...)
	return mockArgs.Get(0).(string), mockArgs.Error(1)
}
