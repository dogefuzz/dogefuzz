package job

import (
	"errors"
	"testing"

	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/test"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionsCheckerJobTestSuite struct {
	suite.Suite

	taskServiceMock           *mocks.TaskServiceMock
	taskInputRequestTopicMock *mocks.TaskInputRequestTopicMock
	env                       Env
}

func TestTransactionsCheckerJobTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionsCheckerJobTestSuite))
}

func (s *TransactionsCheckerJobTestSuite) SetupTest() {
	s.taskServiceMock = new(mocks.TaskServiceMock)
	s.taskInputRequestTopicMock = new(mocks.TaskInputRequestTopicMock)
	s.env = test.NewTestEnv(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, s.taskServiceMock, nil, s.taskInputRequestTopicMock, nil)
}

func (s *TransactionsCheckerJobTestSuite) TestId_ShouldReturnJobName() {
	j := NewTransactionsCheckerJob(s.env)

	id := j.Id()

	assert.Equal(s.T(), "transactions-checker", id)
}

func (s *TransactionsCheckerJobTestSuite) TestCronConfig_ShouldReturnEvery5Seconds() {
	j := NewTransactionsCheckerJob(s.env)

	cronConfig := j.CronConfig()

	assert.Equal(s.T(), "* * * * * *", cronConfig)
}

func (s *TransactionsCheckerJobTestSuite) TestHandler_ShouldFindTransactionsToBeFinishedAndPublishFinishEvent_WhenServiceReturnListOfTransactions() {
	tasks := make([]*dto.TaskDTO, 5)
	for idx := 0; idx < 5; idx++ {
		tasks[idx] = generators.TaskDTOGen()
		s.taskInputRequestTopicMock.On("Publish", bus.TaskInputRequestEvent{TaskId: tasks[idx].Id})
	}
	s.taskServiceMock.On("FindNotFinishedAndHaveDeployedContract").Return(tasks, nil)
	j := NewTransactionsCheckerJob(s.env)

	j.Handler()

	s.taskInputRequestTopicMock.AssertExpectations(s.T())
	s.taskServiceMock.AssertExpectations(s.T())
}

func (s *TransactionsCheckerJobTestSuite) TestHandler_ShouldReturnError_WhenServiceReturnReturnError() {
	tasks := make([]*dto.TaskDTO, 0)
	err := errors.New("error example")
	s.taskServiceMock.On("FindNotFinishedAndHaveDeployedContract").Return(tasks, err)
	j := NewTransactionsCheckerJob(s.env)

	j.Handler()

	s.taskInputRequestTopicMock.AssertNotCalled(s.T(), "Publish")
	s.taskServiceMock.AssertExpectations(s.T())
}

func (s *TransactionsCheckerJobTestSuite) TestHandler_ShouldEmitEvent_WhenServiceReturnAnEmptyListOfTransactions() {
	tasks := make([]*dto.TaskDTO, 0)
	s.taskServiceMock.On("FindNotFinishedAndHaveDeployedContract").Return(tasks, nil)
	j := NewTransactionsCheckerJob(s.env)

	j.Handler()

	s.taskInputRequestTopicMock.AssertNotCalled(s.T(), "Publish")
	s.taskServiceMock.AssertExpectations(s.T())
}
