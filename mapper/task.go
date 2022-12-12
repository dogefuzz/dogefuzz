package mapper

import (
	"strings"
	"time"

	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
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
		Contract:  c.Contract,
		Duration:  int64(c.Duration),
		Detectors: strings.Join(c.Detectors, ";"),
	}
}

func (m *taskMapper) ToDomain(c *dto.TaskDTO) *domain.Task {
	return &domain.Task{
		Id:        c.Id,
		Contract:  c.Contract,
		Duration:  int64(c.Duration),
		Detectors: strings.Join(c.Detectors, ";"),
	}
}

func (m *taskMapper) ToDTO(c *domain.Task) *dto.TaskDTO {
	return &dto.TaskDTO{
		Id:        c.Id,
		Contract:  c.Contract,
		Duration:  time.Duration(c.Duration),
		Detectors: strings.Split(c.Detectors, ";"),
	}
}
