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

type TaskFinishTestSuite struct {
	suite.Suite

	eventBusMock *mocks.EventBusMock
	env          env
}

func TestTaskFinishTestSuite(t *testing.T) {
	suite.Run(t, new(TaskFinishTestSuite))
}

func (s *TaskFinishTestSuite) SetupTest() {
	s.eventBusMock = new(mocks.EventBusMock)
	s.env = test.NewTestEnv(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, s.eventBusMock, nil, nil, nil)
}

func (s *TaskFinishTestSuite) TestPublish_ShouldCallEventBusPublishMethod_WhenReceiveAValidEvent() {
	evt := generators.TaskFinishEventGen()
	s.eventBusMock.On("Publish", TASK_FINISH_TOPIC, []interface{}{evt})
	t := NewTaskFinishTopic(s.env)

	t.Publish(evt)

	s.eventBusMock.AssertExpectations(s.T())
}

func (s *TaskFinishTestSuite) TestSubscribe_ShouldCallEventBusSubscribeMethod_WhenReceiveAValidFunctionHandler() {
	fn := func(foo string) { fmt.Printf("bar %s\n", foo) }
	s.eventBusMock.On("Subscribe", TASK_FINISH_TOPIC, mock.AnythingOfType("func(string)"))
	t := NewTaskFinishTopic(s.env)

	t.Subscribe(fn)

	s.eventBusMock.AssertExpectations(s.T())
}

func (s *TaskFinishTestSuite) TestUnsubscribe_ShouldCallEventBusUnsubscribeMethod_WhenReceiveAValidFunctionHandler() {
	fn := func(foo string) { fmt.Printf("bar %s\n", foo) }
	s.eventBusMock.On("Unsubscribe", TASK_FINISH_TOPIC, mock.AnythingOfType("func(string)"))
	t := NewTaskFinishTopic(s.env)

	t.Unsubscribe(fn)

	s.eventBusMock.AssertExpectations(s.T())
}
