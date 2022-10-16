package service

import (
	"errors"
	"testing"

	"github.com/dogefuzz/dogefuzz/mapper"
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

func (s *ContractServiceTestSuite) TestCreateWithSuccess() {
	newContractDTO := generators.NewContractDTOGen()
	contractDTO := generators.ContractDTOGen()
	contract := generators.ContractGen()
	contractMapperMock := new(mocks.ContractMapperMock)
	contractMapperMock.On("ToDomain", newContractDTO).Return(contract)
	contractMapperMock.On("ToDTO", contract).Return(contractDTO)
	contractRepoMock := new(mocks.ContractRepoMock)
	contractRepoMock.On("Create", contract).Return(nil)

	e := &env{
		contractMapper: contractMapperMock,
		contractRepo:   contractRepoMock,
	}
	contractService := NewContractService(e)
	result, err := contractService.Create(newContractDTO)

	assert.Equal(s.T(), result, contractDTO)
	assert.Nil(s.T(), err)
}

func (s *ContractServiceTestSuite) TestCreateWithFailure() {
	newContractDTO := generators.NewContractDTOGen()
	contractDTO := generators.ContractDTOGen()
	contract := generators.ContractGen()
	contractMapperMock := new(mocks.ContractMapperMock)
	contractMapperMock.On("ToDomain", newContractDTO).Return(contract)
	contractMapperMock.On("ToDTO", contract).Return(contractDTO)
	contractRepoMock := new(mocks.ContractRepoMock)
	expectedError := errors.New("Expected Error")
	contractRepoMock.On("Create", contract).Return(expectedError)

	e := &env{
		contractMapper: contractMapperMock,
		contractRepo:   contractRepoMock,
	}
	contractService := NewContractService(e)
	result, err := contractService.Create(newContractDTO)

	assert.Nil(s.T(), result)
	assert.Equal(s.T(), err, expectedError)
}

type env struct {
	contractMapper mapper.ContractMapper
	contractRepo   repo.ContractRepo
}

func (e *env) ContractMapper() mapper.ContractMapper {
	return e.contractMapper
}

func (e *env) ContractRepo() repo.ContractRepo {
	return e.contractRepo
}
