package mapper

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type TaskMapper interface {
	ToDomainForCreation(c *dto.NewTaskDTO) *entities.Task
	ToDomain(c *dto.TaskDTO) *entities.Task
	ToDTO(c *entities.Task) *dto.TaskDTO
}

type taskMapper struct{}

func NewTaskMapper() *taskMapper {
	return &taskMapper{}
}

func (m *taskMapper) ToDomainForCreation(c *dto.NewTaskDTO) *entities.Task {
	return &entities.Task{
		ContractId: c.ContractId,
		Arguments:  strings.Join(c.Arguments, ";"),
		Expiration: c.Expiration,
		Detectors:  common.JoinOracleTypeList(c.Detectors),
		Status:     c.Status,
	}
}

func (m *taskMapper) ToDomain(c *dto.TaskDTO) *entities.Task {
	return &entities.Task{
		Id:         c.Id,
		ContractId: c.ContractId,
		Arguments:  strings.Join(c.Arguments, ";"),
		Expiration: c.Expiration,
		Detectors:  common.JoinOracleTypeList(c.Detectors),
		Status:     c.Status,
	}
}

func (m *taskMapper) ToDTO(c *entities.Task) *dto.TaskDTO {
	return &dto.TaskDTO{
		Id:         c.Id,
		ContractId: c.ContractId,
		Arguments:  strings.Split(c.Arguments, ";"),
		Expiration: c.Expiration,
		Detectors:  common.SplitOracleTypeString(c.Detectors),
		Status:     c.Status,
	}
}
