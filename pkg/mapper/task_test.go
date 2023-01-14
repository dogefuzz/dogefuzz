package mapper

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TaskMapperTestSuite struct {
	suite.Suite
}

func TestTaskMapperTestSuite(t *testing.T) {
	suite.Run(t, new(TaskMapperTestSuite))
}

func (s *TaskMapperTestSuite) TestMapNewDTOToEntity_ShouldReturnAValidEntity_WhenReveiveAValidNewDTO() {
	newTaskDTO := generators.NewTaskDTOGen()

	m := NewTaskMapper()
	result := m.MapNewDTOToEntity(newTaskDTO)

	expectedResult := entities.Task{
		Arguments:   strings.Join(newTaskDTO.Arguments, ";"),
		StartTime:   newTaskDTO.StartTime,
		Expiration:  newTaskDTO.Expiration,
		Detectors:   common.JoinOracleTypeList(newTaskDTO.Detectors),
		FuzzingType: newTaskDTO.FuzzingType,
		Status:      newTaskDTO.Status,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *TaskMapperTestSuite) TestMapDTOToEntity_ShouldReturnAValidEntity_WhenReveiveAValidDTO() {
	taskDTO := generators.TaskDTOGen()

	m := NewTaskMapper()
	result := m.MapDTOToEntity(taskDTO)

	expectedResult := entities.Task{
		Id:                             taskDTO.Id,
		Arguments:                      strings.Join(taskDTO.Arguments, ";"),
		StartTime:                      taskDTO.StartTime,
		Expiration:                     taskDTO.Expiration,
		Detectors:                      common.JoinOracleTypeList(taskDTO.Detectors),
		FuzzingType:                    taskDTO.FuzzingType,
		AggregatedExecutedInstructions: strings.Join(taskDTO.AggregatedExecutedInstructions, ";"),
		Status:                         taskDTO.Status,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *TaskMapperTestSuite) TestMapEntityToDTO_ShouldReturnAValidDTO_WhenReveiveAValidEntity() {
	entity := generators.TaskGen()

	m := NewTaskMapper()
	result := m.MapEntityToDTO(entity)

	expectedResult := dto.TaskDTO{
		Id:                             entity.Id,
		Arguments:                      strings.Split(entity.Arguments, ";"),
		StartTime:                      entity.StartTime,
		Expiration:                     entity.Expiration,
		Detectors:                      common.SplitOracleTypeString(entity.Detectors),
		FuzzingType:                    entity.FuzzingType,
		AggregatedExecutedInstructions: strings.Split(entity.AggregatedExecutedInstructions, ";"),
		Status:                         entity.Status,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}
