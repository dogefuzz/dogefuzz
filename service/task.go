package service

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskService interface {
	Get(taskId string) (*dto.TaskDTO, error)
	Create(task *dto.NewTaskDTO) (*dto.TaskDTO, error)
	Update(task *dto.TaskDTO) error
	FindNotFinishedTasksThatDontHaveIncompletedTransactions() []*dto.TaskDTO
	FindNotFinishedAndExpired() []*dto.TaskDTO
}

type taskService struct {
	taskRepo   repo.TaskRepo
	taskMapper mapper.TaskMapper
}

func NewTaskService(e Env) *taskService {
	return &taskService{
		taskRepo:   e.TaskRepo(),
		taskMapper: e.TaskMapper(),
	}
}

func (s *taskService) Get(taskId string) (*dto.TaskDTO, error) {
	// TODO: get task
	return nil, nil
}

func (s *taskService) Create(task *dto.NewTaskDTO) (*dto.TaskDTO, error) {
	entity := s.taskMapper.ToDomainForCreation(task)
	err := s.taskRepo.Create(entity)
	if err != nil {
		return nil, err
	}
	return s.taskMapper.ToDTO(entity), nil
}

func (s *taskService) Update(task *dto.TaskDTO) error {
	// TODO: update task in the database
	return nil
}

func (s *taskService) FindNotFinishedTasksThatDontHaveIncompletedTransactions() []*dto.TaskDTO {
	// TODO: Call repository to find tasks that were not finished and don't have incompleted transactions
	return nil
}

func (s *taskService) FindNotFinishedAndExpired() []*dto.TaskDTO {
	// TODO: Call repository to find tasks that were not finished and are expired
	return nil
}
