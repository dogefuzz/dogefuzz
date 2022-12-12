package service

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/common"

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

func (s *ContractServiceTestSuite) TestDeploy_ShouldReturnAddress_WhenProvideNoArgsAndDontOccurFailure() {
	contractService := NewContractService(e)
	contract := generators.CommonContractGen()
	ctx := context.Background()
	expectedAddress := generators.SmartContractGen()
	deployerMock.On("Deploy", ctx, contract).Return(expectedAddress, nil)

	address, err := contractService.Deploy(context.Background(), contract)

	assert.Equal(s.T(), expectedAddress, address)
	assert.Nil(s.T(), err)
}

func (s *ContractServiceTestSuite) TestDeploy_ShouldReturnError_WhenProvideNoArgsAndOccurFailureInDeployment() {
	contractService := NewContractService(e)
	contract := generators.CommonContractGen()
	expectedError := errors.New("expected Error")
	deployerMock.On("Deploy", context.Background(), contract).Return("", expectedError)

	address, err := contractService.Deploy(context.Background(), contract)

	assert.Empty(s.T(), address)
	assert.ErrorIs(s.T(), expectedError, err)
}

func (s *ContractServiceTestSuite) TestDeploy_ShouldReturnAddress_WhenProvideMultipleArgsAndDontOccurFailure() {
	contractService := NewContractService(e)
	args := make([]string, 10)
	gofakeit.Slice(&args)
	contract := generators.CommonContractGen()
	ctx := context.Background()
	expectedAddress := generators.SmartContractGen()
	deployerMock.On("Deploy", packArgsToVariadicFuncParameters(ctx, contract, args)...).Return(expectedAddress, nil)

	address, err := contractService.Deploy(ctx, contract, args...)

	assert.Equal(s.T(), expectedAddress, address)
	assert.Nil(s.T(), err)
}

func (s *ContractServiceTestSuite) TestDeploy_ShouldReturnError_WhenProvideMultipleArgsAndOccurFailureInDeployment() {
	contractService := NewContractService(e)
	args := make([]string, 10)
	gofakeit.Slice(&args)
	contract := generators.CommonContractGen()
	ctx := context.Background()
	expectedError := errors.New("expected Error")
	deployerMock.On("Deploy", packArgsToVariadicFuncParameters(ctx, contract, args)...).Return("", expectedError)

	address, err := contractService.Deploy(ctx, contract, args...)

	assert.Empty(s.T(), address)
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
	contractMapper    mapper.ContractMapper
	transactionMapper mapper.TransactionMapper
	taskMapper        mapper.TaskMapper
	oracleMapper      mapper.OracleMapper
	taskOracleRepo    repo.TaskOracleRepo
	taskRepo          repo.TaskRepo
	contractRepo      repo.ContractRepo
	transactionRepo   repo.TransactionRepo
	oracleRepo        repo.OracleRepo
	deployer          geth.Deployer
}

func (e *env) ContractMapper() mapper.ContractMapper {
	return e.contractMapper
}

func (e *env) TransactionMapper() mapper.TransactionMapper {
	return e.transactionMapper
}

func (e *env) TaskMapper() mapper.TaskMapper {
	return e.taskMapper
}

func (e *env) OracleMapper() mapper.OracleMapper {
	return e.oracleMapper
}

func (e *env) TaskOracleRepo() repo.TaskOracleRepo {
	return e.taskOracleRepo
}

func (e *env) TaskRepo() repo.TaskRepo {
	return e.taskRepo
}

func (e *env) TransactionRepo() repo.TransactionRepo {
	return e.transactionRepo
}

func (e *env) ContractRepo() repo.ContractRepo {
	return e.contractRepo
}

func (e *env) OracleRepo() repo.OracleRepo {
	return e.oracleRepo
}

func (e *env) Deployer() geth.Deployer {
	return e.deployer
}

func packArgsToVariadicFuncParameters(ctx context.Context, contract *common.Contract, args []string) []interface{} {
	parameters := make([]interface{}, len(args)+2)
	parameters[0] = ctx
	parameters[1] = contract
	for idx := 0; idx < len(args); idx++ {
		parameters[idx+2] = args[idx]
	}
	return parameters
}
