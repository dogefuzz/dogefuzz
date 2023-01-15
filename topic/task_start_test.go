package topic

import (
	"fmt"
	"testing"

	"github.com/dogefuzz/dogefuzz/test"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskStartTestSuite struct {
	suite.Suite

	eventBusMock *mocks.EventBusMock
	env          env
}

func TestTaskStartTestSuite(t *testing.T) {
	suite.Run(t, new(TaskStartTestSuite))
}

func (s *TaskStartTestSuite) SetupTest() {
	s.eventBusMock = new(mocks.EventBusMock)
	s.env = test.NewTestEnv(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, s.eventBusMock, nil, nil, nil)
}

func (s *TaskStartTestSuite) TestPublish_ShouldCallEventBusPublishMethod_WhenReceiveAValidEvent() {
	evt := generators.TaskStartEventGen()
	s.eventBusMock.On("Publish", TASK_START_TOPIC, []interface{}{evt})
	t := NewTaskStartTopic(s.env)

	t.Publish(evt)

	s.eventBusMock.AssertExpectations(s.T())
}

func (s *TaskStartTestSuite) TestSubscribe_ShouldCallEventBusSubscribeMethod_WhenReceiveAValidFunctionHandler() {
	fn := func(foo string) { fmt.Printf("bar %s\n", foo) }
	s.eventBusMock.On("Subscribe", TASK_START_TOPIC, mock.AnythingOfType("func(string)"))
	t := NewTaskStartTopic(s.env)

	t.Subscribe(fn)

	s.eventBusMock.AssertExpectations(s.T())
}

func (s *TaskStartTestSuite) TestUnsubscribe_ShouldCallEventBusUnsubscribeMethod_WhenReceiveAValidFunctionHandler() {
	fn := func(foo string) { fmt.Printf("bar %s\n", foo) }
	s.eventBusMock.On("Unsubscribe", TASK_START_TOPIC, mock.AnythingOfType("func(string)"))
	t := NewTaskStartTopic(s.env)

	t.Unsubscribe(fn)

	s.eventBusMock.AssertExpectations(s.T())
}
