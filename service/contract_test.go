package service

import (
	"errors"
	"testing"

	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ContractServiceTestSuite struct {
	suite.Suite
}

func TestContractServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ContractServiceTestSuite))
}

func (s *ContractServiceTestSuite) TestCreate_ShouldReturnCreatedContract_WhenDontOccurFailureDuringCreation() {
	contractService := NewContractService(e)
	newContractDTO := generators.NewContractDTOGen()
	contractDTO := generators.ContractDTOGen()
	contract := generators.ContractGen()
	contractMapperMock.On("ToDomain", newContractDTO).Return(contract)
	contractMapperMock.On("ToDTO", contract).Return(contractDTO)
	contractRepoMock.On("Create", contract).Return(nil)

	result, err := contractService.Create(newContractDTO)

	assert.Equal(s.T(), result, contractDTO)
	assert.Nil(s.T(), err)
}

func (s *ContractServiceTestSuite) TestCreate_ShouldReturnError_WhenOccurFailure() {
	contractService := NewContractService(e)
	newContractDTO := generators.NewContractDTOGen()
	contractDTO := generators.ContractDTOGen()
	contract := generators.ContractGen()
	expectedError := errors.New("expected Error")
	contractMapperMock.On("ToDomain", newContractDTO).Return(contract)
	contractMapperMock.On("ToDTO", contract).Return(contractDTO)
	contractRepoMock.On("Create", contract).Return(expectedError)

	result, err := contractService.Create(newContractDTO)

	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), expectedError, err)
}
