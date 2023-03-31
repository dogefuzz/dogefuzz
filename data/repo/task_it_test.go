package repo

import (
	"reflect"
	"testing"
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/it"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TaskRepoIntegrationTestSuite struct {
	suite.Suite

	env             *it.Env
	repo            interfaces.TaskRepo
	transactionRepo interfaces.TransactionRepo
	contractRepo    interfaces.ContractRepo
}

func TestTaskRepoIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepoIntegrationTestSuite))
}

func (s *TaskRepoIntegrationTestSuite) SetupSuite() {
	config := it.CONFIG
	config.DatabaseName = config.DatabaseName + "_" + uuid.NewString()
	s.env = it.NewEnv(&config)
	s.env.DbConnection().Migrate()
	s.repo = NewTaskRepo(s.env)
	s.transactionRepo = NewTransactionRepo(s.env)
	s.contractRepo = NewContractRepo(s.env)
}

func (s *TaskRepoIntegrationTestSuite) TearDownSuite() {
	s.env.Destroy()
}

func (s *TaskRepoIntegrationTestSuite) TestTaskCreationAndGetMethod_WhenCreatingAValidTask_ShouldNotReturnError() {
	task := generators.TaskGen()

	err := s.repo.Create(s.env.DbConnection().GetDB(), task)
	assert.Nil(s.T(), err)

	foundTask, err := s.repo.Get(s.env.DbConnection().GetDB(), task.Id)
	assert.Nil(s.T(), err)
	assert.True(s.T(), reflect.DeepEqual(task, foundTask))
}

func (s *TaskRepoIntegrationTestSuite) TestTaskUpdate_WhenCreatingAValidTask_ShouldNotReturnError() {
	task := generators.TaskGen()

	err := s.repo.Create(s.env.DbConnection().GetDB(), task)
	assert.Nil(s.T(), err)

	updatedTask := generators.TaskGen()
	updatedTask.Id = task.Id
	err = s.repo.Update(s.env.DbConnection().GetDB(), updatedTask)
	assert.Nil(s.T(), err)

	foundTask, err := s.repo.Get(s.env.DbConnection().GetDB(), task.Id)
	assert.Nil(s.T(), err)
	assert.True(s.T(), reflect.DeepEqual(updatedTask, foundTask))
}

func (s *TaskRepoIntegrationTestSuite) TestFindNotFinishedTasksThatDontHaveIncompletedTransactions_WhenOnlyExistsTasksWithIncompletedTransactions_ShouldReturnEmpty() {
	task1 := generators.TaskGen()
	task1.Status = common.TASK_RUNNING
	err := s.repo.Create(s.env.DbConnection().GetDB(), task1)
	transaction1 := generators.TransactionGen()
	transaction1.Status = common.TRANSACTION_RUNNING
	transaction1.TaskId = task1.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction1)
	assert.Nil(s.T(), err)
	transaction2 := generators.TransactionGen()
	transaction2.Status = common.TRANSACTION_DONE
	transaction2.TaskId = task1.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction2)
	assert.Nil(s.T(), err)

	task2 := generators.TaskGen()
	task2.Status = common.TASK_RUNNING
	err = s.repo.Create(s.env.DbConnection().GetDB(), task2)
	assert.Nil(s.T(), err)
	transaction3 := generators.TransactionGen()
	transaction3.Status = common.TRANSACTION_RUNNING
	transaction3.TaskId = task2.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction3)
	assert.Nil(s.T(), err)
	transaction4 := generators.TransactionGen()
	transaction4.Status = common.TRANSACTION_DONE
	transaction4.TaskId = task2.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction4)
	assert.Nil(s.T(), err)

	foundTasks, err := s.repo.FindNotFinishedTasksThatDontHaveIncompletedTransactions(s.env.DbConnection().GetDB())
	assert.Nil(s.T(), err)
	assert.NotContains(s.T(), foundTasks, task1)
	assert.NotContains(s.T(), foundTasks, task2)
}

func (s *TaskRepoIntegrationTestSuite) TestFindNotFinishedTasksThatDontHaveIncompletedTransactions_WhenThereIsNoRunningTask_ShouldReturnEmpty() {
	task1 := generators.TaskGen()
	task1.Status = common.TASK_DONE
	err := s.repo.Create(s.env.DbConnection().GetDB(), task1)
	assert.Nil(s.T(), err)

	task2 := generators.TaskGen()
	task2.Status = common.TASK_DONE
	err = s.repo.Create(s.env.DbConnection().GetDB(), task2)
	assert.Nil(s.T(), err)

	foundTasks, err := s.repo.FindNotFinishedTasksThatDontHaveIncompletedTransactions(s.env.DbConnection().GetDB())
	assert.Nil(s.T(), err)
	assert.NotContains(s.T(), foundTasks, *task1)
	assert.NotContains(s.T(), foundTasks, *task2)
}

func (s *TaskRepoIntegrationTestSuite) TestFindNotFinishedTasksThatDontHaveIncompletedTransactions_WhenHaveAtLeastOneTaskWithCompletedTransactions_ShouldReturnTheListOfTasksThatSatisfyConstraint() {
	task1 := generators.TaskGen()
	task1.Status = common.TASK_RUNNING
	err := s.repo.Create(s.env.DbConnection().GetDB(), task1)
	transaction1 := generators.TransactionGen()
	transaction1.Status = common.TRANSACTION_RUNNING
	transaction1.TaskId = task1.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction1)
	assert.Nil(s.T(), err)
	transaction2 := generators.TransactionGen()
	transaction2.Status = common.TRANSACTION_DONE
	transaction2.TaskId = task1.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction2)
	assert.Nil(s.T(), err)

	task2 := generators.TaskGen()
	task2.Status = common.TASK_RUNNING
	err = s.repo.Create(s.env.DbConnection().GetDB(), task2)
	assert.Nil(s.T(), err)
	transaction3 := generators.TransactionGen()
	transaction3.Status = common.TRANSACTION_SEND_ERROR
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction3)
	assert.Nil(s.T(), err)
	transaction4 := generators.TransactionGen()
	transaction4.Status = common.TRANSACTION_DONE
	transaction4.TaskId = task1.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction4)
	assert.Nil(s.T(), err)

	foundTasks, err := s.repo.FindNotFinishedTasksThatDontHaveIncompletedTransactions(s.env.DbConnection().GetDB())
	assert.Nil(s.T(), err)
	assert.Contains(s.T(), foundTasks, *task2)
}

func (s *TaskRepoIntegrationTestSuite) TestFindNotFinishedAndExpired_WhenThereIsNoRunningTask_ShouldReturnEmpty() {
	task1 := generators.TaskGen()
	task1.Status = common.TASK_DONE
	err := s.repo.Create(s.env.DbConnection().GetDB(), task1)
	assert.Nil(s.T(), err)

	task2 := generators.TaskGen()
	task2.Status = common.TASK_DONE
	err = s.repo.Create(s.env.DbConnection().GetDB(), task2)
	assert.Nil(s.T(), err)

	foundTasks, err := s.repo.FindNotFinishedAndExpired(s.env.DbConnection().GetDB())
	assert.Nil(s.T(), err)
	assert.NotContains(s.T(), foundTasks, *task1)
	assert.NotContains(s.T(), foundTasks, *task2)
}

func (s *TaskRepoIntegrationTestSuite) TestFindNotFinishedAndExpired_WhenThereIsAtLeastOneRunningTaskButNotExpired_ShouldReturnEmpty() {
	task1 := generators.TaskGen()
	task1.Status = common.TASK_RUNNING
	task1.Expiration = common.Now().Add(1 * time.Hour)
	err := s.repo.Create(s.env.DbConnection().GetDB(), task1)
	assert.Nil(s.T(), err)

	task2 := generators.TaskGen()
	task2.Status = common.TASK_DONE
	task2.Expiration = common.Now().Add(1 * time.Hour)
	err = s.repo.Create(s.env.DbConnection().GetDB(), task2)
	assert.Nil(s.T(), err)

	foundTasks, err := s.repo.FindNotFinishedAndExpired(s.env.DbConnection().GetDB())
	assert.Nil(s.T(), err)
	assert.NotContains(s.T(), foundTasks, *task1)
	assert.NotContains(s.T(), foundTasks, *task2)
}

func (s *TaskRepoIntegrationTestSuite) TestFindNotFinishedAndExpired_WhenThereIsAtLeastOneRunningTaskAndExpired_ShouldTheListOfTasksWithTheseConstraints() {
	task1 := generators.TaskGen()
	task1.Status = common.TASK_RUNNING
	task1.Expiration = common.Now().Add(-1 * time.Hour)
	err := s.repo.Create(s.env.DbConnection().GetDB(), task1)
	assert.Nil(s.T(), err)

	task2 := generators.TaskGen()
	task2.Status = common.TASK_DONE
	task2.Expiration = common.Now().Add(1 * time.Hour)
	err = s.repo.Create(s.env.DbConnection().GetDB(), task2)
	assert.Nil(s.T(), err)

	foundTasks, err := s.repo.FindNotFinishedAndExpired(s.env.DbConnection().GetDB())
	assert.Nil(s.T(), err)
	assert.Contains(s.T(), foundTasks, *task1)
}

func (s *TaskRepoIntegrationTestSuite) TestFindNotFinishedThatHaveDeployedContractAndLimitedPendingTransactions_WhenHavePendingTransactionsLessThanLimit_ShouldReturnTheListOfTasksThatSatisfyConstraint() {
	task1 := generators.TaskGen()
	task1.Status = common.TASK_RUNNING
	err := s.repo.Create(s.env.DbConnection().GetDB(), task1)
	assert.Nil(s.T(), err)
	transaction1 := generators.TransactionGen()
	transaction1.Status = common.TRANSACTION_RUNNING
	transaction1.TaskId = task1.Id
	err = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction1)
	assert.Nil(s.T(), err)
	transaction2 := generators.TransactionGen()
	transaction2.Status = common.TRANSACTION_RUNNING
	transaction2.TaskId = task1.Id
	err = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction2)
	assert.Nil(s.T(), err)
	contract1 := generators.ContractGen()
	contract1.TaskId = task1.Id
	contract1.Status = common.CONTRACT_DEPLOYED
	err = s.contractRepo.Create(s.env.DbConnection().GetDB(), contract1)
	assert.Nil(s.T(), err)

	task2 := generators.TaskGen()
	task2.Status = common.TASK_RUNNING
	err = s.repo.Create(s.env.DbConnection().GetDB(), task2)
	assert.Nil(s.T(), err)
	transaction3 := generators.TransactionGen()
	transaction3.Status = common.TRANSACTION_RUNNING
	transaction3.TaskId = task2.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction3)
	assert.Nil(s.T(), err)
	transaction4 := generators.TransactionGen()
	transaction4.Status = common.TRANSACTION_DONE
	transaction4.TaskId = task2.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction4)
	assert.Nil(s.T(), err)
	contract2 := generators.ContractGen()
	contract2.TaskId = task2.Id
	contract2.Status = common.CONTRACT_DEPLOYED
	err = s.contractRepo.Create(s.env.DbConnection().GetDB(), contract2)
	assert.Nil(s.T(), err)

	foundTasks, err := s.repo.FindNotFinishedThatHaveDeployedContractAndLimitedPendingTransactions(s.env.DbConnection().GetDB(), 2)
	assert.Nil(s.T(), err)
	assert.Contains(s.T(), foundTasks, *task1)
	assert.Contains(s.T(), foundTasks, *task2)
}

func (s *TaskRepoIntegrationTestSuite) TestFindNotFinishedThatHaveDeployedContractAndLimitedPendingTransactions_WhenHavePendingTransactionsGreaterThanLimit_ShouldReturnTheListOfTasksThatSatisfyConstraint() {
	task1 := generators.TaskGen()
	task1.Status = common.TASK_RUNNING
	err := s.repo.Create(s.env.DbConnection().GetDB(), task1)
	assert.Nil(s.T(), err)
	transaction1 := generators.TransactionGen()
	transaction1.Status = common.TRANSACTION_RUNNING
	transaction1.TaskId = task1.Id
	err = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction1)
	assert.Nil(s.T(), err)
	transaction2 := generators.TransactionGen()
	transaction2.Status = common.TRANSACTION_RUNNING
	transaction2.TaskId = task1.Id
	err = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction2)
	assert.Nil(s.T(), err)
	contract1 := generators.ContractGen()
	contract1.TaskId = task1.Id
	contract1.Status = common.CONTRACT_DEPLOYED
	err = s.contractRepo.Create(s.env.DbConnection().GetDB(), contract1)
	assert.Nil(s.T(), err)

	task2 := generators.TaskGen()
	task2.Status = common.TASK_RUNNING
	err = s.repo.Create(s.env.DbConnection().GetDB(), task2)
	assert.Nil(s.T(), err)
	transaction3 := generators.TransactionGen()
	transaction3.Status = common.TRANSACTION_RUNNING
	transaction3.TaskId = task2.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction3)
	assert.Nil(s.T(), err)
	transaction4 := generators.TransactionGen()
	transaction4.Status = common.TRANSACTION_DONE
	transaction4.TaskId = task2.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction4)
	assert.Nil(s.T(), err)
	contract2 := generators.ContractGen()
	contract2.TaskId = task2.Id
	contract2.Status = common.CONTRACT_DEPLOYED
	err = s.contractRepo.Create(s.env.DbConnection().GetDB(), contract2)
	assert.Nil(s.T(), err)

	foundTasks, err := s.repo.FindNotFinishedThatHaveDeployedContractAndLimitedPendingTransactions(s.env.DbConnection().GetDB(), 1)
	assert.Nil(s.T(), err)
	assert.NotContains(s.T(), foundTasks, *task1)
	assert.Contains(s.T(), foundTasks, *task2)
}

func (s *TaskRepoIntegrationTestSuite) TestFindNotFinishedThatHaveDeployedContractAndLimitedPendingTransactions_WhenHaveDeployedContract_ShouldReturnTheListOfTasksThatSatisfyConstraint() {
	task1 := generators.TaskGen()
	task1.Status = common.TASK_RUNNING
	err := s.repo.Create(s.env.DbConnection().GetDB(), task1)
	assert.Nil(s.T(), err)
	transaction1 := generators.TransactionGen()
	transaction1.Status = common.TRANSACTION_RUNNING
	transaction1.TaskId = task1.Id
	err = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction1)
	assert.Nil(s.T(), err)
	transaction2 := generators.TransactionGen()
	transaction2.Status = common.TRANSACTION_RUNNING
	transaction2.TaskId = task1.Id
	err = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction2)
	assert.Nil(s.T(), err)
	contract1 := generators.ContractGen()
	contract1.TaskId = task1.Id
	contract1.Status = common.CONTRACT_DEPLOYED
	err = s.contractRepo.Create(s.env.DbConnection().GetDB(), contract1)
	assert.Nil(s.T(), err)

	task2 := generators.TaskGen()
	task2.Status = common.TASK_RUNNING
	err = s.repo.Create(s.env.DbConnection().GetDB(), task2)
	assert.Nil(s.T(), err)
	transaction3 := generators.TransactionGen()
	transaction3.Status = common.TRANSACTION_RUNNING
	transaction3.TaskId = task2.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction3)
	assert.Nil(s.T(), err)
	transaction4 := generators.TransactionGen()
	transaction4.Status = common.TRANSACTION_DONE
	transaction4.TaskId = task2.Id
	_ = s.transactionRepo.Create(s.env.DbConnection().GetDB(), transaction4)
	assert.Nil(s.T(), err)
	contract2 := generators.ContractGen()
	contract2.TaskId = task2.Id
	contract2.Status = common.CONTRACT_CREATED
	err = s.contractRepo.Create(s.env.DbConnection().GetDB(), contract2)
	assert.Nil(s.T(), err)

	foundTasks, err := s.repo.FindNotFinishedThatHaveDeployedContractAndLimitedPendingTransactions(s.env.DbConnection().GetDB(), 2)
	assert.Nil(s.T(), err)
	assert.Contains(s.T(), foundTasks, *task1)
	assert.NotContains(s.T(), foundTasks, *task2)
}
