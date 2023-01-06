package service

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/data"
	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskService interface {
	Get(taskId string) (*dto.TaskDTO, error)
	Create(task *dto.NewTaskDTO) (*dto.TaskDTO, error)
	Update(task *dto.TaskDTO) error
	FindNotFinishedTasksThatDontHaveIncompletedTransactions() ([]*dto.TaskDTO, error)
	FindNotFinishedAndExpired() ([]*dto.TaskDTO, error)
}

type taskService struct {
	taskRepo   repo.TaskRepo
	taskMapper mapper.TaskMapper
	connection data.Connection
}

func NewTaskService(e Env) *taskService {
	return &taskService{
		taskRepo:   e.TaskRepo(),
		taskMapper: e.TaskMapper(),
		connection: e.DbConnection(),
	}
}

func (s *taskService) Get(taskId string) (*dto.TaskDTO, error) {
	task, err := s.taskRepo.Get(s.connection.GetDB(), taskId)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return s.taskMapper.MapEntityToDTO(task), nil
}

func (s *taskService) Create(task *dto.NewTaskDTO) (*dto.TaskDTO, error) {
	entity := s.taskMapper.MapNewDTOToEntity(task)
	err := s.taskRepo.Create(s.connection.GetDB(), entity)
	if err != nil {
		return nil, err
	}
	return s.taskMapper.MapEntityToDTO(entity), nil
}

func (s *taskService) Update(task *dto.TaskDTO) error {
	entity := s.taskMapper.MapDTOToEntity(task)
	err := s.taskRepo.Update(s.connection.GetDB(), entity)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}

func (s *taskService) FindNotFinishedTasksThatDontHaveIncompletedTransactions() ([]*dto.TaskDTO, error) {
	tasks, err := s.taskRepo.FindNotFinishedTasksThatDontHaveIncompletedTransactions(s.connection.GetDB())
	if err != nil {
		return nil, err
	}
	dtos := make([]*dto.TaskDTO, len(tasks))
	for idx, task := range tasks {
		dtos[idx] = s.taskMapper.MapEntityToDTO(&task)
	}
	return dtos, nil
}

func (s *taskService) FindNotFinishedAndExpired() ([]*dto.TaskDTO, error) {
	tasks, err := s.taskRepo.FindNotFinishedAndExpired(s.connection.GetDB())
	if err != nil {
		return nil, err
	}
	dtos := make([]*dto.TaskDTO, len(tasks))
	for idx, task := range tasks {
		dtos[idx] = s.taskMapper.MapEntityToDTO(&task)
	}
	return dtos, nil
}
