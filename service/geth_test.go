package service

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GethServiceTestSuite struct {
	suite.Suite
}

func TestGethServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GethServiceTestSuite))
}

func (s *GethServiceTestSuite) TestDeploy_ShouldReturnAddress_WhenProvideNoArgsAndDontOccurFailure() {
	contractService := NewGethService(e)
	contract := generators.CommonContractGen()
	ctx := context.Background()
	expectedAddress := generators.SmartContractGen()
	deployerMock.On("Deploy", ctx, contract).Return(expectedAddress, nil)

	address, err := contractService.Deploy(context.Background(), contract)

	assert.Equal(s.T(), expectedAddress, address)
	assert.Nil(s.T(), err)
}

func (s *GethServiceTestSuite) TestDeploy_ShouldReturnError_WhenProvideNoArgsAndOccurFailureInDeployment() {
	contractService := NewGethService(e)
	contract := generators.CommonContractGen()
	expectedError := errors.New("expected Error")
	deployerMock.On("Deploy", context.Background(), contract).Return("", expectedError)

	address, err := contractService.Deploy(context.Background(), contract)

	assert.Empty(s.T(), address)
	assert.ErrorIs(s.T(), expectedError, err)
}

func (s *GethServiceTestSuite) TestDeploy_ShouldReturnAddress_WhenProvideMultipleArgsAndDontOccurFailure() {
	contractService := NewGethService(e)
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

func (s *GethServiceTestSuite) TestDeploy_ShouldReturnError_WhenProvideMultipleArgsAndOccurFailureInDeployment() {
	contractService := NewGethService(e)
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
