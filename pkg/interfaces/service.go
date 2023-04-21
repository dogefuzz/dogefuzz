package interfaces

import (
	"context"
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ContractService interface {
	Find(contractId string) (*dto.ContractDTO, error)
	Get(contractId string) (*dto.ContractDTO, error)
	FindByTaskId(taskId string) (*dto.ContractDTO, error)
	Create(ctr *dto.NewContractDTO) (*dto.ContractDTO, error)
	CreateWithId(ctr *dto.NewContractWithIdDTO) (*dto.ContractDTO, error)
	Update(ctr *dto.ContractDTO) error
}

type FunctionService interface {
	Get(functionId string) (*dto.FunctionDTO, error)
	Create(task *dto.NewFunctionDTO) (*dto.FunctionDTO, error)
	FindByContractId(contractId string) ([]*dto.FunctionDTO, error)
	FindConstructorByContractId(contractId string) (*dto.FunctionDTO, error)
}

type GethService interface {
	Deploy(ctx context.Context, contract *common.Contract, args ...interface{}) (string, string, error)
	BatchCall(ctx context.Context, contract *common.Contract, function *dto.FunctionDTO, inputsByTransactionId map[string][]interface{}) (map[string]string, map[string]error)
}

type ReporterService interface {
	SendReport(ctx context.Context, report common.TaskReport) error
}

type TaskService interface {
	Get(taskId string) (*dto.TaskDTO, error)
	Create(task *dto.NewTaskDTO) (*dto.TaskDTO, error)
	Update(task *dto.TaskDTO) error
	FindNotFinishedTasksThatDontHaveIncompletedTransactions() ([]*dto.TaskDTO, error)
	FindNotFinishedAndExpired() ([]*dto.TaskDTO, error)
	FindNotFinishedAndHaveDeployedContract() ([]*dto.TaskDTO, error)
	FindNotFinishedThatHaveDeployedContractAndLimitedPendingTransactions(limit int) ([]*dto.TaskDTO, error)
}

type TransactionService interface {
	Get(transactionId string) (*dto.TransactionDTO, error)
	Create(transaction *dto.NewTransactionDTO) (*dto.TransactionDTO, error)
	Update(transaction *dto.TransactionDTO) error
	BulkCreate(newTransactions []*dto.NewTransactionDTO) ([]*dto.TransactionDTO, error)
	BulkUpdate(updatedTransactions []*dto.TransactionDTO) error
	FindByHash(hash string) (*dto.TransactionDTO, error)
	FindByTaskId(taskId string) ([]*dto.TransactionDTO, error)
	FindDoneByTaskId(taskId string) ([]*dto.TransactionDTO, error)
	FindDoneTransactionsByFunctionIdAndOrderByTimestamp(functionId string, limit int64) ([]*dto.TransactionDTO, error)
	FindRunningAndCreatedBeforeThreshold(dateThreshold time.Time) ([]*dto.TransactionDTO, error)
}

type VandalService interface {
	GetCFG(ctx context.Context, contract *common.Contract) (*common.CFG, error)
}

type SolidityService interface {
	GetTypeHandlerWithContext(typ abi.Type) (TypeHandler, error)
}
