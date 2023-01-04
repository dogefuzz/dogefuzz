package mocks

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/stretchr/testify/mock"
)

type AgentMock struct {
	mock.Mock
}

func (m *AgentMock) Send(ctx context.Context, contract *common.Contract, functionName string, args ...interface{}) (string, error) {
	arguments := make([]interface{}, 0)
	arguments = append(arguments, ctx)
	arguments = append(arguments, contract)
	arguments = append(arguments, functionName)
	if len(args) > 0 {
		arguments = append(arguments, args...)
	}
	mockArgs := m.Called(arguments...)
	return mockArgs.Get(0).(string), mockArgs.Error(1)
}
