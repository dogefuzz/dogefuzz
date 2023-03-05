package interfaces

import (
	"time"

	"github.com/dogefuzz/dogefuzz/entities"
	"gorm.io/gorm"
)

type ContractRepo interface {
	Create(tx *gorm.DB, contract *entities.Contract) error
	Update(tx *gorm.DB, contract *entities.Contract) error
	FindAll(tx *gorm.DB) ([]entities.Contract, error)
	Find(tx *gorm.DB, id string) (*entities.Contract, error)
	FindByTaskId(tx *gorm.DB, taskId string) (*entities.Contract, error)
}

type FunctionRepo interface {
	Get(tx *gorm.DB, id string) (*entities.Function, error)
	Create(tx *gorm.DB, function *entities.Function) error
	FindByContractId(tx *gorm.DB, contractId string) ([]entities.Function, error)
	FindConstructorByContractId(tx *gorm.DB, contractId string) (*entities.Function, error)
}

type TaskRepo interface {
	Get(tx *gorm.DB, id string) (*entities.Task, error)
	Create(tx *gorm.DB, task *entities.Task) error
	Update(tx *gorm.DB, task *entities.Task) error
	FindNotFinishedTasksThatDontHaveIncompletedTransactions(tx *gorm.DB) ([]entities.Task, error)
	FindNotFinishedAndExpired(tx *gorm.DB) ([]entities.Task, error)
	FindNotFinishedAndHaveDeployedContract(tx *gorm.DB) ([]entities.Task, error)
}

type TransactionRepo interface {
	Get(tx *gorm.DB, id string) (*entities.Transaction, error)
	Create(tx *gorm.DB, transaction *entities.Transaction) error
	Update(tx *gorm.DB, transaction *entities.Transaction) error
	FindByBlockchainHash(tx *gorm.DB, blockchainHash string) (*entities.Transaction, error)
	FindByTaskId(tx *gorm.DB, taskId string) ([]entities.Transaction, error)
	FindDoneByTaskId(tx *gorm.DB, taskId string) ([]entities.Transaction, error)
	FindTransactionsByFunctionNameAndOrderByTimestamp(tx *gorm.DB, functionName string, limit int64) ([]entities.Transaction, error)
	FindRunningAndCreatedBeforeThreshold(tx *gorm.DB, dateThreshold time.Time) ([]entities.Transaction, error)
}
