package repository

import (
	"database/sql"
	"errors"

	"github.com/gongbell/contractfuzzer/db"
	"github.com/gongbell/contractfuzzer/db/domain"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type TransactionRepository interface {
	Create(transaction *domain.Transaction) error
	Update(transaction *domain.Transaction) error
	Find(id string) (*domain.Transaction, error)
	FindByBlockchainHash(blockchainHash string) (*domain.Transaction, error)
	FindTransactionsByTaskId(taskId string) ([]domain.Transaction, error)
	Delete(id string) error
}

type TransactionSQLiteRepository struct {
	manager db.Manager
}

func (r TransactionSQLiteRepository) Init(manager db.Manager) (TransactionSQLiteRepository, error) {
	r.manager = manager
	return r, nil
}

func (r TransactionSQLiteRepository) Create(transaction *domain.Transaction) error {
	transaction.Id = uuid.NewString()
	_, err := r.manager.GetDB().Exec("INSERT INTO transactions(id) values(?)", transaction.Id)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return ErrDuplicate
			}
		}
		return err
	}
	return nil
}

func (r TransactionSQLiteRepository) Update(transaction *domain.Transaction) error {
	query := `
		UPDATE transactions 
		SET 
			blockchain_hash = $2,
			task_id = $3,
			contract_id = $4,
			detected_weaknesses = $5
		WHERE id = $1
	`
	_, err := r.manager.GetDB().Exec(
		query,
		transaction.Id,
		transaction.BlockchainHash,
		transaction.TaskId,
		transaction.ContractId,
		transaction.DetectedWeaknesses,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotExists
			}
		}
		return err
	}
	return nil
}

func (r TransactionSQLiteRepository) Find(id string) (*domain.Transaction, error) {
	row := r.manager.GetDB().QueryRow("SELECT * FROM transactions WHERE id = ?", id)

	var transaction domain.Transaction
	if err := row.Scan(
		&transaction.Id,
		&transaction.BlockchainHash,
		&transaction.TaskId,
		&transaction.ContractId,
		&transaction.DetectedWeaknesses,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &transaction, nil
}

func (r TransactionSQLiteRepository) FindByBlockchainHash(blockchainHash string) (*domain.Transaction, error) {
	row := r.manager.GetDB().QueryRow("SELECT * FROM transactions WHERE blockchain_hash = ?", blockchainHash)

	var transaction domain.Transaction
	if err := row.Scan(
		&transaction.Id,
		&transaction.BlockchainHash,
		&transaction.TaskId,
		&transaction.ContractId,
		&transaction.DetectedWeaknesses,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &transaction, nil
}

func (r TransactionSQLiteRepository) FindTransactionsByTaskId(taskId string) ([]domain.Transaction, error) {
	rows, err := r.manager.GetDB().Query("SELECT * FROM transactions WHERE task_id = ?", taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		if err := rows.Scan(
			&transaction.Id,
			&transaction.BlockchainHash,
			&transaction.TaskId,
			&transaction.ContractId,
			&transaction.DetectedWeaknesses,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r TransactionSQLiteRepository) Delete(id string) error {
	res, err := r.manager.GetDB().Exec("DELETE FROM transactions WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
