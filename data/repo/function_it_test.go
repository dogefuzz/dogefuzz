package repo

import (
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/it"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FunctionRepoIntegrationTestSuite struct {
	suite.Suite

	env  *it.Env
	repo FunctionRepo
}

func TestFunctionRepoIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FunctionRepoIntegrationTestSuite))
}

func (s *FunctionRepoIntegrationTestSuite) SetupSuite() {
	config := it.CONFIG
	config.DatabaseName = config.DatabaseName + "_" + uuid.NewString()
	s.env = it.NewEnv(&config)
	s.env.DbConnection().Migrate()
	s.repo = NewFunctionRepo(s.env)
}

func (s *FunctionRepoIntegrationTestSuite) TearDownSuite() {
	s.env.Destroy()
}

func (s *FunctionRepoIntegrationTestSuite) TestFunctionCreationAndGet_WhenCreatingAValidFunction_ShouldNotReturnError() {
	function := generators.FunctionGen()

	err := s.repo.Create(s.env.DbConnection().GetDB(), function)
	assert.Nil(s.T(), err)

	foundFunction, err := s.repo.Get(s.env.DbConnection().GetDB(), function.Id)
	assert.Nil(s.T(), err)
	assert.True(s.T(), reflect.DeepEqual(function, foundFunction))
}

func (s *FunctionRepoIntegrationTestSuite) TestFindByContractId_WhenExistsMoreThanOneFunctionFromTheSameContract_ShouldReturnListOfFunctions() {
	contractId := gofakeit.LetterN(255)
	function1 := generators.FunctionGen()
	function1.ContractId = contractId
	function1.IsConstructor = false
	err := s.repo.Create(s.env.DbConnection().GetDB(), function1)
	assert.Nil(s.T(), err)
	function2 := generators.FunctionGen()
	function2.ContractId = contractId
	function2.IsConstructor = true
	err = s.repo.Create(s.env.DbConnection().GetDB(), function2)
	assert.Nil(s.T(), err)
	function3 := generators.FunctionGen()
	function3.ContractId = contractId
	function3.IsConstructor = false
	err = s.repo.Create(s.env.DbConnection().GetDB(), function3)
	assert.Nil(s.T(), err)

	functions, err := s.repo.FindByContractId(s.env.DbConnection().GetDB(), contractId)
	assert.Nil(s.T(), err)

	expectedFunctions := []entities.Function{*function1, *function2, *function3}
	assert.ElementsMatch(s.T(), expectedFunctions, functions)
}

func (s *FunctionRepoIntegrationTestSuite) TestFindConstructorByContractId_WhenExistsMultipleFunctionsAndOneOfThemIsAConstructor_ShouldTheValidFunction() {
	contractId := gofakeit.LetterN(255)
	function1 := generators.FunctionGen()
	function1.ContractId = contractId
	function1.IsConstructor = false
	err := s.repo.Create(s.env.DbConnection().GetDB(), function1)
	assert.Nil(s.T(), err)
	function2 := generators.FunctionGen()
	function2.ContractId = contractId
	function2.IsConstructor = true
	err = s.repo.Create(s.env.DbConnection().GetDB(), function2)
	assert.Nil(s.T(), err)
	function3 := generators.FunctionGen()
	function3.ContractId = contractId
	function3.IsConstructor = false
	err = s.repo.Create(s.env.DbConnection().GetDB(), function3)
	assert.Nil(s.T(), err)

	constructor, err := s.repo.FindConstructorByContractId(s.env.DbConnection().GetDB(), contractId)
	assert.Nil(s.T(), err)

	assert.True(s.T(), reflect.DeepEqual(function2, constructor))
}
