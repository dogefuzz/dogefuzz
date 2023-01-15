package mocks

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/stretchr/testify/mock"
)

type TaskInputRequestTopicMock struct {
	mock.Mock
}

func (m *TaskInputRequestTopicMock) Subscribe(fn interface{}) {
	m.Called(fn)
}

func (m *TaskInputRequestTopicMock) Unsubscribe(fn interface{}) {
	m.Called(fn)
}

func (m *TaskInputRequestTopicMock) Publish(e bus.TaskInputRequestEvent) {
	m.Called(e)
}
