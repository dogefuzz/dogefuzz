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

type TasksCheckerJobTestSuite struct {
	suite.Suite

	taskServiceMock     *mocks.TaskServiceMock
	taskFinishTopicMock *mocks.TaskFinishTopicMock
	env                 Env
}

func TestTasksCheckerJobTestSuite(t *testing.T) {
	suite.Run(t, new(TasksCheckerJobTestSuite))
}

func (s *TasksCheckerJobTestSuite) SetupTest() {
	s.taskServiceMock = new(mocks.TaskServiceMock)
	s.taskFinishTopicMock = new(mocks.TaskFinishTopicMock)
	s.env = test.NewTestEnv(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, s.taskServiceMock, s.taskFinishTopicMock, nil, nil)
}

func (s *TasksCheckerJobTestSuite) TestId_ShouldReturnJobName() {
	j := NewTasksCheckerJob(s.env)

	id := j.Id()

	assert.Equal(s.T(), "tasks-checker", id)
}

func (s *TasksCheckerJobTestSuite) TestCronConfig_ShouldReturnEvery5Seconds() {
	j := NewTasksCheckerJob(s.env)

	cronConfig := j.CronConfig()

	assert.Equal(s.T(), "*/5 * * * * *", cronConfig)
}

func (s *TasksCheckerJobTestSuite) TestHandler_ShouldFindTasksToBeFinishedAndPublishFinishEvent_WhenServiceReturnListOfTasks() {
	tasks := make([]*dto.TaskDTO, 5)
	for idx := 0; idx < len(tasks); idx++ {
		tasks[idx] = generators.TaskDTOGen()
		s.taskFinishTopicMock.On("Publish", bus.TaskFinishEvent{TaskId: tasks[idx].Id})
		s.taskServiceMock.On("Update", tasks[idx]).Return(nil)
	}
	s.taskServiceMock.On("FindNotFinishedAndExpired").Return(tasks, nil)
	j := NewTasksCheckerJob(s.env)

	j.Handler()

	s.taskFinishTopicMock.AssertExpectations(s.T())
	s.taskServiceMock.AssertExpectations(s.T())
}

func (s *TasksCheckerJobTestSuite) TestHandler_ShouldReturnError_WhenServiceReturnReturnError() {
	tasks := make([]*dto.TaskDTO, 0)
	err := errors.New("error example")
	s.taskServiceMock.On("FindNotFinishedAndExpired").Return(tasks, err)
	j := NewTasksCheckerJob(s.env)

	j.Handler()

	s.taskFinishTopicMock.AssertNotCalled(s.T(), "Publish")
	s.taskServiceMock.AssertExpectations(s.T())
}

func (s *TasksCheckerJobTestSuite) TestHandler_ShouldEmitEvent_WhenServiceReturnAnEmptyListOfTasks() {
	tasks := make([]*dto.TaskDTO, 0)
	s.taskServiceMock.On("FindNotFinishedAndExpired").Return(tasks, nil)
	j := NewTasksCheckerJob(s.env)

	j.Handler()

	s.taskFinishTopicMock.AssertNotCalled(s.T(), "Publish")
	s.taskServiceMock.AssertExpectations(s.T())
}
