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

type InstrumentExecutionTestSuite struct {
	suite.Suite

	eventBusMock *mocks.EventBusMock
	env          env
}

func TestInstrumentExecutionTestSuite(t *testing.T) {
	suite.Run(t, new(InstrumentExecutionTestSuite))
}

func (s *InstrumentExecutionTestSuite) SetupTest() {
	s.eventBusMock = new(mocks.EventBusMock)
	s.env = test.NewTestEnv(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, s.eventBusMock, nil, nil, nil)
}

func (s *InstrumentExecutionTestSuite) TestPublish_ShouldCallEventBusPublishMethod_WhenReceiveAValidEvent() {
	evt := generators.InstrumentExecutionEventGen()
	s.eventBusMock.On("Publish", INSTRUMENT_EXECUTION_TOPIC, []interface{}{evt})
	t := NewInstrumentExecutionTopic(s.env)

	t.Publish(evt)

	s.eventBusMock.AssertExpectations(s.T())
}

func (s *InstrumentExecutionTestSuite) TestSubscribe_ShouldCallEventBusSubscribeMethod_WhenReceiveAValidFunctionHandler() {
	fn := func(foo string) { fmt.Printf("bar %s\n", foo) }
	s.eventBusMock.On("Subscribe", INSTRUMENT_EXECUTION_TOPIC, mock.AnythingOfType("func(string)"))
	t := NewInstrumentExecutionTopic(s.env)

	t.Subscribe(fn)

	s.eventBusMock.AssertExpectations(s.T())
}

func (s *InstrumentExecutionTestSuite) TestUnsubscribe_ShouldCallEventBusUnsubscribeMethod_WhenReceiveAValidFunctionHandler() {
	fn := func(foo string) { fmt.Printf("bar %s\n", foo) }
	s.eventBusMock.On("Unsubscribe", INSTRUMENT_EXECUTION_TOPIC, mock.AnythingOfType("func(string)"))
	t := NewInstrumentExecutionTopic(s.env)

	t.Unsubscribe(fn)

	s.eventBusMock.AssertExpectations(s.T())
}
