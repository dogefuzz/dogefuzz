package mapper

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionMapperTestSuite struct {
	suite.Suite
}

func TestTransactionMapperTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionMapperTestSuite))
}

func (s *TransactionMapperTestSuite) TestMapNewDTOToEntity_ShouldReturnAValidEntity_WhenReveiveAValidNewDTO() {
	newTransactionDTO := generators.NewTransactionDTOGen()

	m := NewTransactionMapper()
	result := m.MapNewDTOToEntity(newTransactionDTO)

	expectedResult := entities.Transaction{
		Timestamp:  newTransactionDTO.Timestamp,
		TaskId:     newTransactionDTO.TaskId,
		FunctionId: newTransactionDTO.FunctionId,
		Inputs:     strings.Join(newTransactionDTO.Inputs, ";"),
		Status:     newTransactionDTO.Status,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *TransactionMapperTestSuite) TestMapDTOToEntity_ShouldReturnAValidEntity_WhenReveiveAValidDTO() {
	transactionDTO := generators.TransactionDTOGen()

	m := NewTransactionMapper()
	result := m.MapDTOToEntity(transactionDTO)

	coverage := transactionDTO.Coverage
	deltaCoverage := transactionDTO.DeltaCoverage
	deltaMinDistance := transactionDTO.DeltaMinDistance
	criticalInstructionsHits := transactionDTO.CriticalInstructionsHits

	expectedResult := entities.Transaction{
		Id:                       transactionDTO.Id,
		Timestamp:                transactionDTO.Timestamp,
		BlockchainHash:           transactionDTO.BlockchainHash,
		TaskId:                   transactionDTO.TaskId,
		FunctionId:               transactionDTO.FunctionId,
		Inputs:                   strings.Join(transactionDTO.Inputs, ";"),
		DetectedWeaknesses:       strings.Join(transactionDTO.DetectedWeaknesses, ";"),
		ExecutedInstructions:     strings.Join(transactionDTO.ExecutedInstructions, ";"),
		Coverage:                 coverage,
		DeltaCoverage:            deltaCoverage,
		DeltaMinDistance:         deltaMinDistance,
		Status:                   transactionDTO.Status,
		CriticalInstructionsHits: criticalInstructionsHits,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *TransactionMapperTestSuite) TestMapEntityToDTO_ShouldReturnAValidDTO_WhenReveiveAValidEntity() {
	entity := generators.TransactionGen()

	m := NewTransactionMapper()
	result := m.MapEntityToDTO(entity)

	expectedResult := dto.TransactionDTO{
		Id:                       entity.Id,
		Timestamp:                entity.Timestamp,
		BlockchainHash:           entity.BlockchainHash,
		TaskId:                   entity.TaskId,
		FunctionId:               entity.FunctionId,
		Inputs:                   strings.Split(entity.Inputs, ";"),
		DetectedWeaknesses:       strings.Split(entity.DetectedWeaknesses, ";"),
		ExecutedInstructions:     strings.Split(entity.ExecutedInstructions, ";"),
		Coverage:                 entity.Coverage,
		DeltaCoverage:            entity.DeltaCoverage,
		DeltaMinDistance:         entity.DeltaMinDistance,
		Status:                   entity.Status,
		CriticalInstructionsHits: entity.CriticalInstructionsHits,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}
