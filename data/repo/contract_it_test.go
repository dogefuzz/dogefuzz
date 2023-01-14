package repo

import (
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/it"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ContractRepoIntegrationTestSuite struct {
	suite.Suite

	env  *it.Env
	repo interfaces.ContractRepo
}

func TestContractRepoIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ContractRepoIntegrationTestSuite))
}

func (s *ContractRepoIntegrationTestSuite) SetupSuite() {
	config := it.CONFIG
	config.DatabaseName = config.DatabaseName + "_" + uuid.NewString()
	s.env = it.NewEnv(&config)
	s.env.DbConnection().Migrate()
	s.repo = NewContractRepo(s.env)
}

func (s *ContractRepoIntegrationTestSuite) TearDownSuite() {
	s.env.Destroy()
}

func (s *ContractRepoIntegrationTestSuite) TestContractCreationAndFindMethods_WhenCreatingAValidContract_ShouldNotReturnError() {
	contract := generators.ContractGen()

	err := s.repo.Create(s.env.DbConnection().GetDB(), contract)
	assert.Nil(s.T(), err)

	foundContract, err := s.repo.Find(s.env.DbConnection().GetDB(), contract.Id)
	assert.Nil(s.T(), err)
	assert.True(s.T(), reflect.DeepEqual(contract, foundContract))

	foundContractByTaskId, err := s.repo.FindByTaskId(s.env.DbConnection().GetDB(), contract.TaskId)
	assert.Nil(s.T(), err)
	assert.True(s.T(), reflect.DeepEqual(contract, foundContractByTaskId))
}

func (s *ContractRepoIntegrationTestSuite) TestContractUpdate_WhenCreatingAValidContract_ShouldNotReturnError() {
	contract := generators.ContractGen()
	updatedContract := contract
	updatedContract.Address = generators.SmartContractGen()
	updatedContract.Source = gofakeit.LetterN(255)
	updatedContract.Name = gofakeit.Name()

	err := s.repo.Create(s.env.DbConnection().GetDB(), contract)
	assert.Nil(s.T(), err)

	err = s.repo.Update(s.env.DbConnection().GetDB(), updatedContract)
	assert.Nil(s.T(), err)

	foundContract, err := s.repo.Find(s.env.DbConnection().GetDB(), contract.Id)
	assert.Nil(s.T(), err)
	assert.True(s.T(), reflect.DeepEqual(updatedContract, foundContract))
}
