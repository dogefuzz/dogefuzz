package repo

import (
	"database/sql"
	"errors"

	"github.com/dogefuzz/dogefuzz/data"
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type TaskRepo interface {
	Create(task *entities.Task) error
	Find(id string) (*entities.Task, error)
	Delete(id string) error
}

type taskRepo struct {
	connection data.Connection
}

func NewTaskRepo(e Env) *taskRepo {
	return &taskRepo{connection: e.DbConnection()}
}

func (r *taskRepo) Create(task *entities.Task) error {
	task.Id = uuid.NewString()
	_, err := r.connection.GetDB().Exec("INSERT INTO tasks(id) values(?)", task.Id)
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

func (r *taskRepo) Find(id string) (*entities.Task, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	var task entities.Task
	if err := row.Scan(&task.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) Delete(id string) error {
	res, err := r.connection.GetDB().Exec("DELETE FROM tasks WHERE id = ?", id)
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
