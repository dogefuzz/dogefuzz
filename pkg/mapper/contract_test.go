package mapper

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
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

func (s *ContractMapperTestSuite) TestMapNewDTOToEntity_ShouldReturnValidEntity_WhenReceiveAValidNewContractDTO() {
	newContractDTO := generators.NewContractDTOGen()

	m := NewContractMapper()
	result := m.MapNewDTOToEntity(newContractDTO)

	expectedResult := entities.Contract{
		TaskId:        newContractDTO.TaskId,
		Source:        newContractDTO.Source,
		CompiledCode:  newContractDTO.CompiledCode,
		AbiDefinition: newContractDTO.AbiDefinition,
		Name:          newContractDTO.Name,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *ContractMapperTestSuite) TestMapDTOToEntity_ShouldReturnValidEntity_WhenReceiveAValidContractDTO() {
	contractDTO := generators.ContractDTOGen()

	m := NewContractMapper()
	result := m.MapDTOToEntity(contractDTO)

	expectedCFG, _ := json.Marshal(contractDTO.CFG)
	expectedDistanceMap, _ := json.Marshal(contractDTO.DistanceMap)
	expectedResult := entities.Contract{
		Id:            contractDTO.Id,
		TaskId:        contractDTO.TaskId,
		Address:       contractDTO.Address,
		Source:        contractDTO.Source,
		CompiledCode:  contractDTO.CompiledCode,
		AbiDefinition: contractDTO.AbiDefinition,
		Name:          contractDTO.Name,
		CFG:           string(expectedCFG),
		DistanceMap:   string(expectedDistanceMap),
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *ContractMapperTestSuite) TestMapEntityToDTO_ShouldReturnValidDTO_WhenReceiveValidEntity() {
	contract := generators.ContractGen()

	m := NewContractMapper()
	result := m.MapEntityToDTO(contract)

	var expectedCFG common.CFG
	_ = json.Unmarshal([]byte(contract.CFG), &expectedCFG)
	var expectedDistanceMap common.DistanceMap
	_ = json.Unmarshal([]byte(contract.DistanceMap), &expectedDistanceMap)
	expectedResult := dto.ContractDTO{
		Id:            contract.Id,
		TaskId:        contract.TaskId,
		Address:       contract.Address,
		Source:        contract.Source,
		CompiledCode:  contract.CompiledCode,
		AbiDefinition: contract.AbiDefinition,
		Name:          contract.Name,
		CFG:           expectedCFG,
		DistanceMap:   expectedDistanceMap,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *ContractMapperTestSuite) TestMapEntityToDTO_ShouldReturnValidDTOWithEmptyCFGAndDistanceMap_WhenReceiveValidEntityWithInvalidCFGAndDistanceMap() {
	contract := generators.ContractGen()
	contract.CFG = gofakeit.Word()
	contract.DistanceMap = gofakeit.Word()

	m := NewContractMapper()
	result := m.MapEntityToDTO(contract)

	expectedResult := dto.ContractDTO{
		Id:            contract.Id,
		TaskId:        contract.TaskId,
		Address:       contract.Address,
		Source:        contract.Source,
		CompiledCode:  contract.CompiledCode,
		AbiDefinition: contract.AbiDefinition,
		Name:          contract.Name,
		CFG:           common.CFG{},
		DistanceMap:   common.DistanceMap(nil),
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *ContractMapperTestSuite) TestMapDTOToCommon_ShouldReturnValidCommonContract_WhenReceiveValidDTO() {
	contractDTO := generators.ContractDTOGen()

	m := NewContractMapper()
	result := m.MapDTOToCommon(contractDTO)

	expectedResult := common.Contract{
		Address:       contractDTO.Address,
		CompiledCode:  contractDTO.CompiledCode,
		AbiDefinition: contractDTO.AbiDefinition,
		Name:          contractDTO.Name,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}
