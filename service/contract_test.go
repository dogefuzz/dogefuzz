package service

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/repo"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/mocks"
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
	newContractDTO := generators.NewContractDTOGen()
	contractDTO := generators.ContractDTOGen()
	contract := generators.ContractGen()

	contractMapperMock.On("ToDomain", newContractDTO).Return(contract)
	contractMapperMock.On("ToDTO", contract).Return(contractDTO)
	contractRepoMock.On("Create", contract).Return(nil)

	contractService := NewContractService(e)
	result, err := contractService.Create(newContractDTO)

	assert.Equal(s.T(), result, contractDTO)
	assert.Nil(s.T(), err)
}

func (s *ContractServiceTestSuite) TestCreate_ShouldReturnErro_WhenOccurFailure() {
	newContractDTO := generators.NewContractDTOGen()
	contractDTO := generators.ContractDTOGen()
	contract := generators.ContractGen()
	expectedError := errors.New("Expected Error")

	contractMapperMock.On("ToDomain", newContractDTO).Return(contract)
	contractMapperMock.On("ToDTO", contract).Return(contractDTO)
	contractRepoMock.On("Create", contract).Return(expectedError)

	contractService := NewContractService(e)
	result, err := contractService.Create(newContractDTO)

	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), expectedError, err)
}

func (s *ContractServiceTestSuite) TestDeploy_ShouldReturnAddress_WhenDontOccurFailure() {
	expectedAddress := gofakeit.HexUint256()
	contract := generators.CommonContractGen()
	deployerMock.On("Deploy").Return(expectedAddress, nil)

	args := make([]string, 0)

	contractService := NewContractService(e)
	address, err := contractService.Deploy(context.Background(), contract, args...)

	assert.Equal(s.T(), expectedAddress, address)
	assert.Nil(s.T(), err)
}

func (s *ContractServiceTestSuite) TestDeploy_ShouldReturnError_WhenOccurFailureInDeployment() {
	expectedError := errors.New("Expected Error")
	contract := generators.CommonContractGen()
	deployerMock.On("Deploy").Return(nil, expectedError)

	contractService := NewContractService(e)
	address, err := contractService.Deploy(context.Background(), contract)

	assert.Nil(s.T(), address)
	assert.ErrorIs(s.T(), expectedError, err)
}

var contractMapperMock = new(mocks.ContractMapperMock)
var contractRepoMock = new(mocks.ContractRepoMock)
var deployerMock = new(mocks.DeployerMock)
var e = &env{
	contractMapper: contractMapperMock,
	contractRepo:   contractRepoMock,
	deployer:       deployerMock,
}

type env struct {
	contractMapper mapper.ContractMapper
	contractRepo   repo.ContractRepo
	deployer       geth.Deployer
}

func (e *env) ContractMapper() mapper.ContractMapper {
	return e.contractMapper
}

func (e *env) ContractRepo() repo.ContractRepo {
	return e.contractRepo
}

func (e *env) Deployer() geth.Deployer {
	return e.deployer
}
