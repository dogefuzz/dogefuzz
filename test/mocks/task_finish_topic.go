package mocks

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/stretchr/testify/mock"
)

type TaskFinishTopicMock struct {
	mock.Mock
}

func (m *TaskFinishTopicMock) Subscribe(fn interface{}) {
	m.Called(fn)
}

func (m *TaskFinishTopicMock) Unsubscribe(fn interface{}) {
	m.Called(fn)
}

func (m *TaskFinishTopicMock) Publish(e bus.TaskFinishEvent) {
	m.Called(e)
}
