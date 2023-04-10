package repo

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/it"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionRepoIntegrationTestSuite struct {
	suite.Suite

	env          *it.Env
	repo         interfaces.TransactionRepo
	functionRepo interfaces.FunctionRepo
}

func TestTransactionRepoIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepoIntegrationTestSuite))
}

func (s *TransactionRepoIntegrationTestSuite) SetupSuite() {
	config := it.CONFIG
	config.DatabaseName = config.DatabaseName + "_" + uuid.NewString()
	s.env = it.NewEnv(&config)
	s.env.DbConnection().Migrate()
	s.repo = NewTransactionRepo(s.env)
	s.functionRepo = NewFunctionRepo(s.env)
}

func (s *TransactionRepoIntegrationTestSuite) TearDownSuite() {
	s.env.Destroy()
}

func (s *TransactionRepoIntegrationTestSuite) TestTransactionCreationAndGetMethod_WhenCreatingAValidTransaction_ShouldNotReturnError() {
	transaction := generators.TransactionGen()

	err := s.repo.Create(s.env.DbConnection().GetDB(), transaction)
	assert.Nil(s.T(), err)

	foundTransaction, err := s.repo.Get(s.env.DbConnection().GetDB(), transaction.Id)
	assert.Nil(s.T(), err)
	assert.True(s.T(), reflect.DeepEqual(transaction, foundTransaction))

	foundTransactionByBlockchainHash, err := s.repo.FindByBlockchainHash(s.env.DbConnection().GetDB(), transaction.BlockchainHash)
	assert.Nil(s.T(), err)
	assert.True(s.T(), reflect.DeepEqual(transaction, foundTransactionByBlockchainHash))
}

func (s *TransactionRepoIntegrationTestSuite) TestTransactionUpdate_WhenCreatingAValidTransaction_ShouldNotReturnError() {
	transaction := generators.TransactionGen()

	err := s.repo.Create(s.env.DbConnection().GetDB(), transaction)
	assert.Nil(s.T(), err)

	updatedTransaction := generators.TransactionGen()
	updatedTransaction.Id = transaction.Id
	err = s.repo.Update(s.env.DbConnection().GetDB(), updatedTransaction)
	assert.Nil(s.T(), err)

	foundTransaction, err := s.repo.Get(s.env.DbConnection().GetDB(), transaction.Id)
	assert.Nil(s.T(), err)
	assert.True(s.T(), reflect.DeepEqual(updatedTransaction, foundTransaction))
}

func (s *TransactionRepoIntegrationTestSuite) TestFindByTaskId_WhenExistsMultipleTransactionsWithTheSameTaskId_ShouldReturnListOfTransactions() {
	taskId := gofakeit.LetterN(255)
	transaction1 := generators.TransactionGen()
	transaction1.TaskId = taskId
	err := s.repo.Create(s.env.DbConnection().GetDB(), transaction1)
	assert.Nil(s.T(), err)
	transaction2 := generators.TransactionGen()
	transaction2.TaskId = taskId
	err = s.repo.Create(s.env.DbConnection().GetDB(), transaction2)
	assert.Nil(s.T(), err)
	transaction3 := generators.TransactionGen()
	transaction3.TaskId = taskId
	err = s.repo.Create(s.env.DbConnection().GetDB(), transaction3)
	assert.Nil(s.T(), err)

	transactions, err := s.repo.FindByTaskId(s.env.DbConnection().GetDB(), taskId)
	assert.Nil(s.T(), err)
	assert.ElementsMatch(s.T(), []entities.Transaction{*transaction1, *transaction2, *transaction3}, transactions)
}

func (s *TransactionRepoIntegrationTestSuite) TestFindDoneTransactionsByFunctionNameAndOrderByTimestamp_WhenExistsMultipleTransactionsWithTheSameFunction_ShouldReturnListOfTransactions() {
	function := generators.FunctionGen()
	err := s.functionRepo.Create(s.env.DbConnection().GetDB(), function)
	assert.Nil(s.T(), err)
	transaction1 := generators.TransactionGen()
	transaction1.Status = common.TRANSACTION_DONE
	transaction1.FunctionId = function.Id
	transaction1.Timestamp = common.Now()
	err = s.repo.Create(s.env.DbConnection().GetDB(), transaction1)
	assert.Nil(s.T(), err)
	transaction2 := generators.TransactionGen()
	transaction2.Status = common.TRANSACTION_DONE
	transaction2.FunctionId = function.Id
	transaction2.Timestamp = common.Now().Add(1 * time.Hour)
	err = s.repo.Create(s.env.DbConnection().GetDB(), transaction2)
	assert.Nil(s.T(), err)
	transaction3 := generators.TransactionGen()
	transaction3.Status = common.TRANSACTION_RUNNING
	transaction3.FunctionId = function.Id
	transaction3.Timestamp = common.Now().Add(2 * time.Hour)
	err = s.repo.Create(s.env.DbConnection().GetDB(), transaction3)
	assert.Nil(s.T(), err)

	transactions, err := s.repo.FindDoneTransactionsByFunctionNameAndOrderByTimestamp(s.env.DbConnection().GetDB(), function.Name, 10)
	assert.Nil(s.T(), err)
	assert.True(s.T(), sort.SliceIsSorted(transactions, func(i, j int) bool {
		return transactions[i].Timestamp.Before(transactions[j].Timestamp)
	}))
	assert.ElementsMatch(s.T(), []entities.Transaction{*transaction1, *transaction2}, transactions)
}
