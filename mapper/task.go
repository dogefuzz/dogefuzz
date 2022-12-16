package mapper

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type TaskMapper interface {
	ToDomainForCreation(c *dto.NewTaskDTO) *domain.Task
	ToDomain(c *dto.TaskDTO) *domain.Task
	ToDTO(c *domain.Task) *dto.TaskDTO
}

type taskMapper struct{}

func NewTaskMapper() *taskMapper {
	return &taskMapper{}
}

func (m *taskMapper) ToDomainForCreation(c *dto.NewTaskDTO) *domain.Task {
	return &domain.Task{
		ContractId: c.ContractId,
		Arguments:  strings.Join(c.Arguments, ";"),
		Expiration: c.Expiration,
		Detectors:  common.JoinOracleTypeList(c.Detectors),
		Status:     c.Status,
	}
}

func (m *taskMapper) ToDomain(c *dto.TaskDTO) *domain.Task {
	return &domain.Task{
		Id:         c.Id,
		ContractId: c.ContractId,
		Arguments:  strings.Join(c.Arguments, ";"),
		Expiration: c.Expiration,
		Detectors:  common.JoinOracleTypeList(c.Detectors),
		Status:     c.Status,
	}
}

func (m *taskMapper) ToDTO(c *domain.Task) *dto.TaskDTO {
	return &dto.TaskDTO{
		Id:         c.Id,
		ContractId: c.ContractId,
		Arguments:  strings.Split(c.Arguments, ";"),
		Expiration: c.Expiration,
		Detectors:  common.SplitOracleTypeString(c.Detectors),
		Status:     c.Status,
	}
}
