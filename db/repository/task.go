package repository

import (
	"database/sql"
	"errors"

	"github.com/gongbell/contractfuzzer/db"
	"github.com/gongbell/contractfuzzer/db/domain"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	Find(id string) (*domain.Task, error)
	Delete(id string) error
}

type TaskSQLiteRepository struct {
	manager db.Manager
}

func (r TaskSQLiteRepository) Init(manager db.Manager) (TaskSQLiteRepository, error) {
	r.manager = manager
	return r, nil
}

func (r TaskSQLiteRepository) Create(task *domain.Task) error {
	task.Id = uuid.NewString()
	_, err := r.manager.GetDB().Exec("INSERT INTO tasks(id) values(?)", task.Id)
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

func (r TaskSQLiteRepository) Find(id string) (*domain.Task, error) {
	row := r.manager.GetDB().QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	var task domain.Task
	if err := row.Scan(&task.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &task, nil
}

func (r TaskSQLiteRepository) Delete(id string) error {
	res, err := r.manager.GetDB().Exec("DELETE FROM tasks WHERE id = ?", id)
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
