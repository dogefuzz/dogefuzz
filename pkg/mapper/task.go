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
		Arguments:      strings.Join(c.Arguments, ";"),
		Duration:       c.Duration,
		StartTime:      c.StartTime,
		DeploymentTime: c.DeploymentTime,
		Expiration:     c.Expiration,
		Detectors:      common.JoinOracleTypeList(c.Detectors),
		FuzzingType:    c.FuzzingType,
		Status:         c.Status,
	}
}

func (m *taskMapper) MapDTOToEntity(c *dto.TaskDTO) *entities.Task {
	return &entities.Task{
		Id:                             c.Id,
		Arguments:                      strings.Join(c.Arguments, ";"),
		Duration:                       c.Duration,
		StartTime:                      c.StartTime,
		DeploymentTime:                 c.DeploymentTime,
		Expiration:                     c.Expiration,
		Detectors:                      common.JoinOracleTypeList(c.Detectors),
		FuzzingType:                    c.FuzzingType,
		AggregatedExecutedInstructions: strings.Join(c.AggregatedExecutedInstructions, ";"),
		Status:                         c.Status,
	}
}

func (m *taskMapper) MapEntityToDTO(c *entities.Task) *dto.TaskDTO {
	var arguments []string
	if c.Arguments != "" {
		arguments = strings.Split(c.Arguments, ";")
	}

	var detectors []common.OracleType
	if c.Detectors != "" {
		detectors = common.SplitOracleTypeString(c.Detectors)
	}

	var aggregatedExecutedInstructions []string
	if c.AggregatedExecutedInstructions != "" {
		aggregatedExecutedInstructions = strings.Split(c.AggregatedExecutedInstructions, ";")
	}

	return &dto.TaskDTO{
		Id:                             c.Id,
		Arguments:                      arguments,
		Duration:                       c.Duration,
		StartTime:                      c.StartTime,
		DeploymentTime:                 c.DeploymentTime,
		Expiration:                     c.Expiration,
		Detectors:                      detectors,
		FuzzingType:                    c.FuzzingType,
		AggregatedExecutedInstructions: aggregatedExecutedInstructions,
		Status:                         c.Status,
	}
}
