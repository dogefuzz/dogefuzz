package mapper

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type taskMapper struct{}

func NewTaskMapper() *taskMapper {
	return &taskMapper{}
}

func (m *taskMapper) MapNewDTOToEntity(c *dto.NewTaskDTO) *entities.Task {
	return &entities.Task{
		Arguments:   strings.Join(c.Arguments, ";"),
		StartTime:   c.StartTime,
		Expiration:  c.Expiration,
		Detectors:   common.JoinOracleTypeList(c.Detectors),
		FuzzingType: c.FuzzingType,
		Status:      c.Status,
	}
}

func (m *taskMapper) MapDTOToEntity(c *dto.TaskDTO) *entities.Task {
	return &entities.Task{
		Id:                             c.Id,
		Arguments:                      strings.Join(c.Arguments, ";"),
		StartTime:                      c.StartTime,
		Expiration:                     c.Expiration,
		Detectors:                      common.JoinOracleTypeList(c.Detectors),
		FuzzingType:                    c.FuzzingType,
		AggregatedExecutedInstructions: strings.Join(c.AggregatedExecutedInstructions, ";"),
		Status:                         c.Status,
	}
}

func (m *taskMapper) MapEntityToDTO(c *entities.Task) *dto.TaskDTO {
	return &dto.TaskDTO{
		Id:                             c.Id,
		Arguments:                      strings.Split(c.Arguments, ";"),
		StartTime:                      c.StartTime,
		Expiration:                     c.Expiration,
		Detectors:                      common.SplitOracleTypeString(c.Detectors),
		FuzzingType:                    c.FuzzingType,
		AggregatedExecutedInstructions: strings.Split(c.AggregatedExecutedInstructions, ";"),
		Status:                         c.Status,
	}
}
