package repo

import (
	"errors"
	"time"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionRepo struct {
}

func NewTransactionRepo(e Env) *transactionRepo {
	return &transactionRepo{}
}

func (r *transactionRepo) Get(tx *gorm.DB, id string) (*entities.Transaction, error) {
	var transaction entities.Transaction
	if err := tx.First(&transaction, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepo) Create(tx *gorm.DB, transaction *entities.Transaction) error {
	transaction.Id = uuid.NewString()
	return tx.Create(transaction).Error
}

func (r *transactionRepo) Update(tx *gorm.DB, updatedTransaction *entities.Transaction) error {
	var transaction entities.Transaction
	if err := tx.First(&transaction, "id = ?", updatedTransaction.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotExists
		}
		return err
	}
	return tx.Model(&transaction).Updates(updatedTransaction).Error
}

func (r *transactionRepo) FindByBlockchainHash(tx *gorm.DB, blockchainHash string) (*entities.Transaction, error) {
	var transaction entities.Transaction
	if err := tx.First(&transaction, "blockchain_hash = ?", blockchainHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepo) FindByTaskId(tx *gorm.DB, taskId string) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	if err := tx.Where("task_id = ?", taskId).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepo) FindDoneByTaskId(tx *gorm.DB, taskId string) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	if err := tx.Where("status = ?", common.TASK_DONE).Where("task_id = ?", taskId).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepo) FindDoneTransactionsByFunctionIdAndOrderByTimestamp(tx *gorm.DB, functionId string, limit int64) ([]entities.Transaction, error) {
	var function entities.Function
	if err := tx.Where("id = ?", functionId).Find(&function).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}

	var transactions []entities.Transaction
	if err := tx.Where("function_id", function.Id).Where("status", common.TRANSACTION_DONE).Order("timestamp").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepo) FindRunningAndCreatedBeforeThreshold(tx *gorm.DB, dateThreshold time.Time) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	if err := tx.Where("status = ?", common.TRANSACTION_RUNNING).Where("timestamp < ?", dateThreshold).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepo) FindTimeTakenToWeakness(tx *gorm.DB, taskId string, weaknessType common.OracleType) (uint32, error) {
	var totalTimeInSeconds uint32
	query := `
		SELECT 
			strftime('%s', timestamp) - strftime('%s', (
				SELECT MIN(timestamp) 
				FROM transactions
				WHERE task_id = ?)
				)
		FROM transactions 
		WHERE detected_weaknesses like ?
		ORDER BY timestamp
		LIMIT 1	
	`

	if err := tx.Raw(query, taskId, "%"+weaknessType+"%").Scan(&totalTimeInSeconds).Error; err != nil {
		return 0, err
	}
	return totalTimeInSeconds, nil
}
