package repo

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/entities"
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

func (r *transactionRepo) FindTransactionsByFunctionNameAndOrderByTimestamp(tx *gorm.DB, functionName string, limit int64) ([]entities.Transaction, error) {
	var function entities.Function
	if err := tx.Where("name = ?", functionName).Find(&function).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}

	var transactions []entities.Transaction
	if err := tx.Where("function_id", function.Id).Order("timestamp").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
