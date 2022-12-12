package service

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/repo"
)

var ErrOraclesNotFound = errors.New("no oracles were found associated with the current transaction")

type OracleService interface {
	FindByTaskId(taskId string) ([]*dto.OracleDTO, error)
	FindByName(name string) (*dto.OracleDTO, error)
}

type oracleService struct {
	transactionRepo repo.TransactionRepo
	taskRepo        repo.TaskRepo
	taskOracleRepo  repo.TaskOracleRepo
	oracleRepo      repo.OracleRepo
	oracleMapper    mapper.OracleMapper
}

func NewOracleService(e Env) *oracleService {
	return &oracleService{
		transactionRepo: e.TransactionRepo(),
		taskRepo:        e.TaskRepo(),
		taskOracleRepo:  e.TaskOracleRepo(),
		oracleRepo:      e.OracleRepo(),
		oracleMapper:    e.OracleMapper(),
	}
}

func (s *oracleService) FindByTaskId(taskId string) ([]*dto.OracleDTO, error) {

	task, err := s.taskRepo.Find(taskId)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrOraclesNotFound
		}
		return nil, err
	}

	tasksOracles, err := s.taskOracleRepo.FindByTaskId(task.Id)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrOraclesNotFound
		}
		return nil, err
	}
	var oracleIds []string
	for _, to := range tasksOracles {
		oracleIds = append(oracleIds, to.OracleId)
	}

	oracles, err := s.oracleRepo.FindByIds(oracleIds)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrOraclesNotFound
		}
		return nil, err
	}

	dtos := make([]*dto.OracleDTO, len(oracles))
	for idx, o := range oracles {
		dtos[idx] = s.oracleMapper.ToDTO(&o)
	}
	return dtos, nil
}

func (s *oracleService) FindByName(name string) (*dto.OracleDTO, error) {
	oracle, err := s.oracleRepo.FindByName(name)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrOraclesNotFound
		}
		return nil, err
	}
	return s.oracleMapper.ToDTO(oracle), nil
}
