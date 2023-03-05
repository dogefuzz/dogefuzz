package mocks

import (
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/stretchr/testify/mock"
)

type TaskServiceMock struct {
	mock.Mock
}

func (m *TaskServiceMock) Get(taskId string) (*dto.TaskDTO, error) {
	args := m.Called(taskId)
	return args.Get(0).(*dto.TaskDTO), args.Error(1)
}

func (m *TaskServiceMock) Create(task *dto.NewTaskDTO) (*dto.TaskDTO, error) {
	args := m.Called(task)
	return args.Get(0).(*dto.TaskDTO), args.Error(1)
}

func (m *TaskServiceMock) Update(task *dto.TaskDTO) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *TaskServiceMock) FindNotFinishedTasksThatDontHaveIncompletedTransactions() ([]*dto.TaskDTO, error) {
	args := m.Called()
	return args.Get(0).([]*dto.TaskDTO), args.Error(1)
}

func (m *TaskServiceMock) FindNotFinishedAndExpired() ([]*dto.TaskDTO, error) {
	args := m.Called()
	return args.Get(0).([]*dto.TaskDTO), args.Error(1)
}

func (m *TaskServiceMock) FindNotFinishedAndHaveDeployedContract() ([]*dto.TaskDTO, error) {
	args := m.Called()
	return args.Get(0).([]*dto.TaskDTO), args.Error(1)
}
