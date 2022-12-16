package repo

import (
	"database/sql"
	"errors"

	"github.com/dogefuzz/dogefuzz/db"
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type FunctionRepo interface {
	GetById(id string) (*domain.Function, error)
	Create(function *domain.Function) error
}

type functionRepo struct {
	connection db.Connection
}

func NewFunctionRepo(e Env) *functionRepo {
	return &functionRepo{connection: e.DbConnection()}
}

func (r *functionRepo) GetById(id string) (*domain.Function, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM functions WHERE id = ?", id)

	var function domain.Function
	if err := row.Scan(
		&function.Id,
		&function.Name,
		&function.NumberOfArgs,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &function, nil
}

func (r *functionRepo) Create(function *domain.Function) error {
	function.Id = uuid.NewString()
	_, err := r.connection.GetDB().Exec("INSERT INTO functions(id, name, number_of_args) values(?, ?, ?)", function.Id, function.Name, function.NumberOfArgs)
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
