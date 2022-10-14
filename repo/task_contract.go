package repo

import (
	"database/sql"
	"errors"

	"github.com/dogefuzz/dogefuzz/db"
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/mattn/go-sqlite3"
)

type TaskContractRepo interface {
	Create(taskContract *domain.TaskContract) error
	FindByTaskIdAndContractId(taskId string, contractId string) (*domain.TaskContract, error)
	FindByTaskId(taskId string) ([]domain.TaskContract, error)
	Delete(taskId string, contractId string) error
}

type taskContractRepo struct {
	connection db.Connection
}

func NewTaskContractRepo(e Env) *taskContractRepo {
	return &taskContractRepo{connection: e.DbConnection()}
}

func (r *taskContractRepo) Create(taskContract *domain.TaskContract) error {
	_, err := r.connection.GetDB().Exec("INSERT INTO tasks_contracts(task_id, contract_id) values(?, ?)", taskContract.TaskId, taskContract.ContractId)
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

func (r *taskContractRepo) FindByTaskIdAndContractId(taskId string, contractId string) (*domain.TaskContract, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM tasks_contracts WHERE task_id = ? AND contract_id = ?", taskId, contractId)

	var taskContract domain.TaskContract
	if err := row.Scan(&taskContract.TaskId, &taskContract.ContractId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &taskContract, nil
}

func (r *taskContractRepo) FindByTaskId(taskId string) ([]domain.TaskContract, error) {
	rows, err := r.connection.GetDB().Query("SELECT * FROM tasks_contracts WHERE task_id = ?", taskId)
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

func (r *taskContractRepo) Delete(taskId string, contractId string) error {
	res, err := r.connection.GetDB().Exec("DELETE FROM tasks_contracts WHERE task_id = ? AND contract_id = ?", taskId, contractId)
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
