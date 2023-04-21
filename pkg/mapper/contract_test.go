package mapper

import (
	"encoding/json"
	"reflect"
	"strconv"
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
		TaskId:             newContractDTO.TaskId,
		Status:             newContractDTO.Status,
		Source:             newContractDTO.Source,
		DeploymentBytecode: newContractDTO.DeploymentBytecode,
		RuntimeBytecode:    newContractDTO.RuntimeBytecode,
		AbiDefinition:      newContractDTO.AbiDefinition,
		Name:               newContractDTO.Name,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *ContractMapperTestSuite) TestMapDTOToEntity_ShouldReturnValidEntity_WhenReceiveAValidContractDTO() {
	contractDTO := generators.ContractDTOGen()

	m := NewContractMapper()
	result := m.MapDTOToEntity(contractDTO)

	expectedCFG, _ := json.Marshal(contractDTO.CFG)
	expectedDistanceMap, _ := json.Marshal(contractDTO.DistanceMap)
	targetInstructionsFreq := strconv.FormatUint(contractDTO.TargetInstructionsFreq, 10)

	expectedResult := entities.Contract{
		Id:                     contractDTO.Id,
		TaskId:                 contractDTO.TaskId,
		Status:                 contractDTO.Status,
		Address:                contractDTO.Address,
		Source:                 contractDTO.Source,
		DeploymentBytecode:     contractDTO.DeploymentBytecode,
		RuntimeBytecode:        contractDTO.RuntimeBytecode,
		AbiDefinition:          contractDTO.AbiDefinition,
		Name:                   contractDTO.Name,
		CFG:                    string(expectedCFG),
		DistanceMap:            string(expectedDistanceMap),
		TargetInstructionsFreq: targetInstructionsFreq,
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

	var targetInstructionsFreq uint64
	if contract.TargetInstructionsFreq != "" {
		val, err := strconv.ParseUint(contract.TargetInstructionsFreq, 10, 64)
		if err != nil {
			panic(err)
		}
		targetInstructionsFreq = val
	}

	expectedResult := dto.ContractDTO{
		Id:                     contract.Id,
		TaskId:                 contract.TaskId,
		Status:                 contract.Status,
		Address:                contract.Address,
		Source:                 contract.Source,
		DeploymentBytecode:     contract.DeploymentBytecode,
		RuntimeBytecode:        contract.RuntimeBytecode,
		AbiDefinition:          contract.AbiDefinition,
		Name:                   contract.Name,
		CFG:                    expectedCFG,
		DistanceMap:            expectedDistanceMap,
		TargetInstructionsFreq: targetInstructionsFreq,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *ContractMapperTestSuite) TestMapEntityToDTO_ShouldReturnValidDTOWithEmptyCFGAndDistanceMap_WhenReceiveValidEntityWithInvalidCFGAndDistanceMap() {
	contract := generators.ContractGen()
	contract.CFG = gofakeit.Word()
	contract.DistanceMap = gofakeit.Word()

	m := NewContractMapper()
	result := m.MapEntityToDTO(contract)

	var targetInstructionsFreq uint64
	if contract.TargetInstructionsFreq != "" {
		val, err := strconv.ParseUint(contract.TargetInstructionsFreq, 10, 64)
		if err != nil {
			panic(err)
		}
		targetInstructionsFreq = val
	}

	expectedResult := dto.ContractDTO{
		Id:                     contract.Id,
		TaskId:                 contract.TaskId,
		Status:                 contract.Status,
		Address:                contract.Address,
		Source:                 contract.Source,
		DeploymentBytecode:     contract.DeploymentBytecode,
		RuntimeBytecode:        contract.RuntimeBytecode,
		AbiDefinition:          contract.AbiDefinition,
		Name:                   contract.Name,
		CFG:                    common.CFG{},
		DistanceMap:            common.DistanceMap(nil),
		TargetInstructionsFreq: targetInstructionsFreq,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}

func (s *ContractMapperTestSuite) TestMapDTOToCommon_ShouldReturnValidCommonContract_WhenReceiveValidDTO() {
	contractDTO := generators.ContractDTOGen()

	m := NewContractMapper()
	result := m.MapDTOToCommon(contractDTO)

	expectedResult := common.Contract{
		Address:            contractDTO.Address,
		DeploymentBytecode: contractDTO.DeploymentBytecode,
		RuntimeBytecode:    contractDTO.RuntimeBytecode,
		AbiDefinition:      contractDTO.AbiDefinition,
		Name:               contractDTO.Name,
	}
	assert.True(s.T(), reflect.DeepEqual(expectedResult, *result))
}
