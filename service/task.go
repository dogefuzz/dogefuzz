package service

import (
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/repo"
)

type TaskService interface {
	Create(task *dto.NewTaskDTO) (*dto.TaskDTO, error)
	AddOracle(taskId string, oracleId string) error
	AddContract(taskId string, contractId string) error
}

type taskService struct {
	taskRepo         repo.TaskRepo
	taskOracleRepo   repo.TaskOracleRepo
	taskContractRepo repo.TaskContractRepo
	taskMapper       mapper.TaskMapper
}

func NewTaskService(e Env) *taskService {
	return &taskService{
		taskRepo:       e.TaskRepo(),
		taskOracleRepo: e.TaskOracleRepo(),
		taskMapper:     e.TaskMapper(),
	}
}

func (s *taskService) Create(task *dto.NewTaskDTO) (*dto.TaskDTO, error) {
	entity := s.taskMapper.ToDomainForCreation(task)
	err := s.taskRepo.Create(entity)
	if err != nil {
		return nil, err
	}
	return s.taskMapper.ToDTO(entity), nil
}

func (s *taskService) AddOracle(taskId string, oracleId string) error {
	taskOracle := domain.TaskOracle{}
	taskOracle.OracleId = oracleId
	taskOracle.TaskId = taskId
	return s.taskOracleRepo.Create(&taskOracle)
}

func (s *taskService) AddContract(taskId string, contractId string) error {
	taskContract := domain.TaskContract{}
	taskContract.ContractId = contractId
	taskContract.TaskId = taskId
	return s.taskContractRepo.Create(&taskContract)
}
