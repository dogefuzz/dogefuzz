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

type ContractMapperTestSuite struct {
	suite.Suite
}

func TestContractMapperTestSuite(t *testing.T) {
	suite.Run(t, new(ContractMapperTestSuite))
}

func (s *ContractMapperTestSuite) TestToDomain() {
	newContractDTO := generators.NewContractDTOGen()

	m := NewContractMapper()
	result := m.ToDomain(newContractDTO)

	expectedResult := entities.Contract{
		Name:   newContractDTO.Name,
		Source: newContractDTO.Source,
	}
	assert.True(s.T(), reflect.DeepEqual(*result, expectedResult))
}

func (s *ContractMapperTestSuite) TestToDTO() {
	contract := generators.ContractGen()

	m := NewContractMapper()
	result := m.ToDTO(contract)

	expectedResult := dto.ContractDTO{
		Id:      contract.Id,
		Address: contract.Address,
		Source:  contract.Source,
		Name:    contract.Name,
	}
	assert.True(s.T(), reflect.DeepEqual(*result, expectedResult))
}
