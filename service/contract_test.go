package service

import (
	"errors"
	"testing"

	"github.com/dogefuzz/dogefuzz/test"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ContractServiceTestSuite struct {
	suite.Suite

	contractMapperMock *mocks.ContractMapperMock
	contractRepoMock   *mocks.ContractRepoMock
	deployerMock       *mocks.DeployerMock
	connectionMock     *mocks.ConnectionMock
	env                *test.TestEnv
}

func TestContractServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ContractServiceTestSuite))
}

func (s *ContractServiceTestSuite) SetupTest() {
	s.contractMapperMock = new(mocks.ContractMapperMock)
	s.contractRepoMock = new(mocks.ContractRepoMock)
	s.deployerMock = new(mocks.DeployerMock)
	s.connectionMock = new(mocks.ConnectionMock)
	s.env = test.NewTestEnv(s.contractMapperMock, nil, nil, nil, nil, s.contractRepoMock, nil, nil, s.deployerMock, nil, s.connectionMock, nil, nil, nil, nil)
}

func (s *ContractServiceTestSuite) TestCreate_ShouldReturnCreatedContract_WhenDontOccurFailureDuringCreation() {
	dummyGormDb := new(gorm.DB)
	contractService := NewContractService(s.env)
	newContractDTO := generators.NewContractDTOGen()
	contractDTO := generators.ContractDTOGen()
	contract := generators.ContractGen()

	s.contractMapperMock.On("MapNewDTOToEntity", newContractDTO).Return(contract)
	s.contractMapperMock.On("MapEntityToDTO", contract).Return(contractDTO)
	s.contractRepoMock.On("Create", dummyGormDb, contract).Return(nil)
	s.connectionMock.On("GetDB").Return(dummyGormDb)
	result, err := contractService.Create(newContractDTO)

	assert.Equal(s.T(), result, contractDTO)
	assert.Nil(s.T(), err)
}

func (s *ContractServiceTestSuite) TestCreate_ShouldReturnError_WhenOccurFailure() {
	dummyGormDb := new(gorm.DB)
	contractService := NewContractService(s.env)
	newContractDTO := generators.NewContractDTOGen()
	contractDTO := generators.ContractDTOGen()
	contract := generators.ContractGen()
	expectedError := errors.New("expected Error")

	s.contractMapperMock.On("MapNewDTOToEntity", newContractDTO).Return(contract)
	s.contractMapperMock.On("MapEntityToDTO", contract).Return(contractDTO)
	s.contractRepoMock.On("Create", dummyGormDb, contract).Return(expectedError)
	s.connectionMock.On("GetDB").Return(dummyGormDb)

	result, err := contractService.Create(newContractDTO)

	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), expectedError, err)
}
