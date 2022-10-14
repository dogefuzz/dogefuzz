package repository

import (
	"database/sql"
	"errors"

	"github.com/gongbell/contractfuzzer/db"
	"github.com/gongbell/contractfuzzer/db/domain"
	"github.com/mattn/go-sqlite3"
)

type TaskOracleRepository interface {
	Create(oracle *domain.TaskOracle) error
	FindByTaskIdAndOracleId(taskId string, oracleId string) (*domain.TaskOracle, error)
	FindByTaskId(taskId string) ([]domain.TaskOracle, error)
	Delete(taskId string, oracleId string) error
}

type TaskOracleSQLiteRepository struct {
	manager db.Manager
}

func (r TaskOracleSQLiteRepository) Init(manager db.Manager) (TaskOracleSQLiteRepository, error) {
	r.manager = manager
	return r, nil
}

func (r TaskOracleSQLiteRepository) Create(taskOracle *domain.TaskOracle) error {
	_, err := r.manager.GetDB().Exec("INSERT INTO tasks_oracles(task_id, oracle_id) values(?, ?)", taskOracle.TaskId, taskOracle.OracleId)
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

func (r TaskOracleSQLiteRepository) FindByTaskIdAndOracleId(taskId string, oracleId string) (*domain.TaskOracle, error) {
	row := r.manager.GetDB().QueryRow("SELECT * FROM tasks_oracles WHERE task_id = ? AND oracle_id = ?", taskId, oracleId)

	var taskOracle domain.TaskOracle
	if err := row.Scan(&taskOracle.TaskId, &taskOracle.OracleId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &taskOracle, nil
}

func (r TaskOracleSQLiteRepository) FindByTaskId(taskId string) ([]domain.TaskOracle, error) {
	rows, err := r.manager.GetDB().Query("SELECT * FROM tasks_oracles WHERE task_id = ?", taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskOracles []domain.TaskOracle
	for rows.Next() {
		var taskOracle domain.TaskOracle
		if err := rows.Scan(&taskOracle.TaskId, &taskOracle.OracleId); err != nil {
			return nil, err
		}
		taskOracles = append(taskOracles, taskOracle)
	}
	return taskOracles, nil
}

func (r TaskOracleSQLiteRepository) Delete(taskId string, oracleId string) error {
	res, err := r.manager.GetDB().Exec("DELETE FROM tasks_oracles WHERE task_id = ? AND oracle_id = ?", taskId, oracleId)
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
