package mapper

import (
	"reflect"
	"testing"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FunctionMapperTestSuite struct {
	suite.Suite
}

func TestFunctionMapperTestSuite(t *testing.T) {
	suite.Run(t, new(FunctionMapperTestSuite))
}

func (s *FunctionMapperTestSuite) TestMapNewDTOToEntity_ShouldReturnAValidEntity_WhenReveiveAValidNewDTO() {
	newFunctionDTO := generators.NewFunctionDTOGen()

	m := NewFunctionMapper()
	result := m.MapNewDTOToEntity(newFunctionDTO)

	expectedResult := entities.Function{
		Name:         newFunctionDTO.Name,
		NumberOfArgs: newFunctionDTO.NumberOfArgs,
		Callable:     newFunctionDTO.Callable,
		Type:         newFunctionDTO.Type,
		ContractId:   newFunctionDTO.ContractId,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *FunctionMapperTestSuite) TestMapDTOToEntity_ShouldReturnAValidEntity_WhenReveiveAValidDTO() {
	functionDTO := generators.FunctionDTOGen()

	m := NewFunctionMapper()
	result := m.MapDTOToEntity(functionDTO)

	expectedResult := entities.Function{
		Id:           functionDTO.Id,
		Name:         functionDTO.Name,
		NumberOfArgs: functionDTO.NumberOfArgs,
		Callable:     functionDTO.Callable,
		Type:         functionDTO.Type,
		ContractId:   functionDTO.ContractId,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *FunctionMapperTestSuite) TestMapEntityToDTO_ShouldReturnAValidDTO_WhenReveiveAValidEntity() {
	entity := generators.FunctionGen()

	m := NewFunctionMapper()
	result := m.MapEntityToDTO(entity)

	expectedResult := dto.FunctionDTO{
		Name:         entity.Name,
		NumberOfArgs: entity.NumberOfArgs,
		Callable:     entity.Callable,
		Type:         entity.Type,
		ContractId:   entity.ContractId,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}
