package mocks

import "github.com/stretchr/testify/mock"

type EventBusMock struct {
	mock.Mock
}

func (m *EventBusMock) Subscribe(topic string, fn interface{}) {
	m.Called(topic, fn)
}

func (m *EventBusMock) Unsubscribe(topic string, fn interface{}) {
	m.Called(topic, fn)
}

func (m *EventBusMock) SubscribeOnce(topic string, fn interface{}) {
	m.Called(topic, fn)
}

func (m *EventBusMock) Publish(topic string, args ...interface{}) {
	m.Called(topic, args)
}
