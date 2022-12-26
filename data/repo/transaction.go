package repo

import (
	"database/sql"
	"errors"

	"github.com/dogefuzz/dogefuzz/data"
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type TransactionRepo interface {
	Create(transaction *entities.Transaction) error
	Update(transaction *entities.Transaction) error
	Find(id string) (*entities.Transaction, error)
	FindByBlockchainHash(blockchainHash string) (*entities.Transaction, error)
	FindTransactionsByTaskId(taskId string) ([]entities.Transaction, error)
	Delete(id string) error
}

type transactionRepo struct {
	connection data.Connection
}

func NewTransactionRepo(e Env) *transactionRepo {
	return &transactionRepo{connection: e.DbConnection()}
}

func (r *transactionRepo) Create(transaction *entities.Transaction) error {
	transaction.Id = uuid.NewString()
	_, err := r.connection.GetDB().Exec("INSERT INTO transactions(id) values(?)", transaction.Id)
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

func (r *transactionRepo) Update(transaction *entities.Transaction) error {
	query := `
		UPDATE transactions 
		SET 
			blockchain_hash = $2,
			task_id = $3,
			transaction_id = $4,
			detected_weaknesses = $5
		WHERE id = $1
	`
	_, err := r.connection.GetDB().Exec(
		query,
		transaction.Id,
		transaction.BlockchainHash,
		transaction.TaskId,
		transaction.FunctionId,
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

func (r *transactionRepo) Find(id string) (*entities.Transaction, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM transactions WHERE id = ?", id)

	var transaction entities.Transaction
	if err := row.Scan(
		&transaction.Id,
		&transaction.BlockchainHash,
		&transaction.TaskId,
		&transaction.FunctionId,
		&transaction.DetectedWeaknesses,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepo) FindByBlockchainHash(blockchainHash string) (*entities.Transaction, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM transactions WHERE blockchain_hash = ?", blockchainHash)

	var transaction entities.Transaction
	if err := row.Scan(
		&transaction.Id,
		&transaction.BlockchainHash,
		&transaction.TaskId,
		&transaction.FunctionId,
		&transaction.DetectedWeaknesses,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepo) FindTransactionsByTaskId(taskId string) ([]entities.Transaction, error) {
	rows, err := r.connection.GetDB().Query("SELECT * FROM transactions WHERE task_id = ?", taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []entities.Transaction
	for rows.Next() {
		var transaction entities.Transaction
		if err := rows.Scan(
			&transaction.Id,
			&transaction.BlockchainHash,
			&transaction.TaskId,
			&transaction.FunctionId,
			&transaction.DetectedWeaknesses,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepo) Delete(id string) error {
	res, err := r.connection.GetDB().Exec("DELETE FROM transactions WHERE id = ?", id)
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
