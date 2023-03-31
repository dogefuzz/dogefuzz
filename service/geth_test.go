package service

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/dogefuzz/dogefuzz/test"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GethServiceTestSuite struct {
	suite.Suite

	contractMapperMock *mocks.ContractMapperMock
	contractRepoMock   *mocks.ContractRepoMock
	deployerMock       *mocks.DeployerMock
	env                *test.TestEnv
}

func TestGethServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GethServiceTestSuite))
}

func (s *GethServiceTestSuite) SetupTest() {
	s.contractMapperMock = new(mocks.ContractMapperMock)
	s.contractRepoMock = new(mocks.ContractRepoMock)
	s.deployerMock = new(mocks.DeployerMock)
	s.env = test.NewTestEnv(s.contractMapperMock, nil, nil, nil, nil, s.contractRepoMock, nil, nil, s.deployerMock, nil, nil, nil, nil, nil, nil, generators.ConfigGen())
}

func (s *GethServiceTestSuite) TestDeploy_ShouldReturnAddress_WhenProvideNoArgsAndDontOccurFailure() {
	contractService := NewGethService(s.env)
	contract := generators.CommonContractGen()
	ctx := context.Background()
	expectedAddress := generators.SmartContractGen()
	expectedTx := gofakeit.LetterN(20)
	s.deployerMock.On("Deploy", ctx, contract).Return(expectedAddress, expectedTx, nil)

	address, tx, err := contractService.Deploy(context.Background(), contract)

	assert.Equal(s.T(), expectedAddress, address)
	assert.Equal(s.T(), expectedTx, tx)
	assert.Nil(s.T(), err)
}

func (s *GethServiceTestSuite) TestDeploy_ShouldReturnError_WhenProvideNoArgsAndOccurFailureInDeployment() {
	contractService := NewGethService(s.env)
	contract := generators.CommonContractGen()
	expectedError := errors.New("expected Error")
	s.deployerMock.On("Deploy", context.Background(), contract).Return("", "", expectedError)

	address, tx, err := contractService.Deploy(context.Background(), contract)

	assert.Empty(s.T(), address)
	assert.Empty(s.T(), tx)
	assert.ErrorIs(s.T(), expectedError, err)
}

func (s *GethServiceTestSuite) TestDeploy_ShouldReturnAddress_WhenProvideMultipleArgsAndDontOccurFailure() {
	contractService := NewGethService(s.env)
	args := make([]interface{}, 10)
	gofakeit.Slice(&args)
	contract := generators.CommonContractGen()
	ctx := context.Background()
	expectedAddress := generators.SmartContractGen()
	expectedTx := gofakeit.LetterN(20)
	s.deployerMock.On("Deploy", packArgsToVariadicFuncParameters(ctx, contract, args)...).Return(expectedAddress, expectedTx, nil)

	address, tx, err := contractService.Deploy(ctx, contract, args...)

	assert.Equal(s.T(), expectedTx, tx)
	assert.Equal(s.T(), expectedAddress, address)
	assert.Nil(s.T(), err)
}

func (s *GethServiceTestSuite) TestDeploy_ShouldReturnError_WhenProvideMultipleArgsAndOccurFailureInDeployment() {
	contractService := NewGethService(s.env)
	args := make([]interface{}, 10)
	gofakeit.Slice(&args)
	contract := generators.CommonContractGen()
	ctx := context.Background()
	expectedError := errors.New("expected Error")
	s.deployerMock.On("Deploy", packArgsToVariadicFuncParameters(ctx, contract, args)...).Return("", "", expectedError)

	address, tx, err := contractService.Deploy(ctx, contract, args...)

	assert.Empty(s.T(), address)
	assert.Empty(s.T(), tx)
	assert.ErrorIs(s.T(), expectedError, err)
}
