package repository

import (
	"database/sql"
	"errors"

	"github.com/gongbell/contractfuzzer/db"
	"github.com/gongbell/contractfuzzer/db/domain"
	"github.com/mattn/go-sqlite3"
)

type TaskContractRepository interface {
	Create(taskContract *domain.TaskContract) error
	FindByTaskIdAndContractId(taskId string, contractId string) (*domain.TaskContract, error)
	FindByTaskId(taskId string) ([]domain.TaskContract, error)
	Delete(taskId string, contractId string) error
}

type TaskContractSQLiteRepository struct {
	manager db.Manager
}

func (r TaskContractSQLiteRepository) Init(manager db.Manager) (TaskContractSQLiteRepository, error) {
	r.manager = manager
	return r, nil
}

func (r TaskContractSQLiteRepository) Create(taskContract *domain.TaskContract) error {
	_, err := r.manager.GetDB().Exec("INSERT INTO tasks_contracts(task_id, contract_id) values(?, ?)", taskContract.TaskId, taskContract.ContractId)
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

func (r TaskContractSQLiteRepository) FindByTaskIdAndContractId(taskId string, contractId string) (*domain.TaskContract, error) {
	row := r.manager.GetDB().QueryRow("SELECT * FROM tasks_contracts WHERE task_id = ? AND contract_id = ?", taskId, contractId)

	var taskContract domain.TaskContract
	if err := row.Scan(&taskContract.TaskId, &taskContract.ContractId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &taskContract, nil
}

func (r TaskContractSQLiteRepository) FindByTaskId(taskId string) ([]domain.TaskContract, error) {
	rows, err := r.manager.GetDB().Query("SELECT * FROM tasks_contracts WHERE task_id = ?", taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskContracts []domain.TaskContract
	for rows.Next() {
		var taskContract domain.TaskContract
		if err := rows.Scan(&taskContract.TaskId, &taskContract.ContractId); err != nil {
			return nil, err
		}
		taskContracts = append(taskContracts, taskContract)
	}
	return taskContracts, nil
}

func (r TaskContractSQLiteRepository) Delete(taskId string, contractId string) error {
	res, err := r.manager.GetDB().Exec("DELETE FROM tasks_contracts WHERE task_id = ? AND contract_id = ?", taskId, contractId)
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
