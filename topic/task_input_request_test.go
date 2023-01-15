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

type TaskInputRequestTestSuite struct {
	suite.Suite

	eventBusMock *mocks.EventBusMock
	env          env
}

func TestTaskInputRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TaskInputRequestTestSuite))
}

func (s *TaskInputRequestTestSuite) SetupTest() {
	s.eventBusMock = new(mocks.EventBusMock)
	s.env = test.NewTestEnv(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, s.eventBusMock, nil, nil, nil)
}

func (s *TaskInputRequestTestSuite) TestPublish_ShouldCallEventBusPublishMethod_WhenReceiveAValidEvent() {
	evt := generators.TaskInputRequestEventGen()
	s.eventBusMock.On("Publish", TASK_INPUT_REQUEST_TOPIC, []interface{}{evt})
	t := NewTaskInputRequestTopic(s.env)

	t.Publish(evt)

	s.eventBusMock.AssertExpectations(s.T())
}

func (s *TaskInputRequestTestSuite) TestSubscribe_ShouldCallEventBusSubscribeMethod_WhenReceiveAValidFunctionHandler() {
	fn := func(foo string) { fmt.Printf("bar %s\n", foo) }
	s.eventBusMock.On("Subscribe", TASK_INPUT_REQUEST_TOPIC, mock.AnythingOfType("func(string)"))
	t := NewTaskInputRequestTopic(s.env)

	t.Subscribe(fn)

	s.eventBusMock.AssertExpectations(s.T())
}

func (s *TaskInputRequestTestSuite) TestUnsubscribe_ShouldCallEventBusUnsubscribeMethod_WhenReceiveAValidFunctionHandler() {
	fn := func(foo string) { fmt.Printf("bar %s\n", foo) }
	s.eventBusMock.On("Unsubscribe", TASK_INPUT_REQUEST_TOPIC, mock.AnythingOfType("func(string)"))
	t := NewTaskInputRequestTopic(s.env)

	t.Unsubscribe(fn)

	s.eventBusMock.AssertExpectations(s.T())
}
