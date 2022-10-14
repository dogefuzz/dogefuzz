package repo

import (
	"database/sql"
	"errors"

	"github.com/dogefuzz/dogefuzz/db"
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type TaskRepo interface {
	Create(task *domain.Task) error
	Find(id string) (*domain.Task, error)
	Delete(id string) error
}

type taskRepo struct {
	connection db.Connection
}

func NewTaskRepo(e Env) *taskRepo {
	return &taskRepo{connection: e.DbConnection()}
}

func (r *taskRepo) Create(task *domain.Task) error {
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

func (r *taskRepo) Find(id string) (*domain.Task, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	var task domain.Task
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
