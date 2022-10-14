package repo

import (
	"database/sql"
	"errors"

	"github.com/dogefuzz/dogefuzz/db"
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type OracleRepo interface {
	Create(oracle *domain.Oracle) error
	Find(id string) (*domain.Oracle, error)
	FindByName(name string) (*domain.Oracle, error)
	FindByIds(ids []string) ([]domain.Oracle, error)
	FindByNames(name []string) ([]domain.Oracle, error)
	Delete(id string) error
}

type oracleRepo struct {
	connection db.Connection
}

func NewOracleRepo(e Env) *oracleRepo {
	return &oracleRepo{connection: e.DbConnection()}
}

func (r *oracleRepo) Create(oracle *domain.Oracle) error {
	oracle.Id = uuid.NewString()
	_, err := r.connection.GetDB().Exec("INSERT INTO oracles(id, name) values(?, ?)", oracle.Id, oracle.Name)
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

func (r *oracleRepo) Find(id string) (*domain.Oracle, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM oracles WHERE id = ?", id)

	var oracle domain.Oracle
	if err := row.Scan(&oracle.Id, &oracle.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &oracle, nil
}

func (r *oracleRepo) FindByName(name string) (*domain.Oracle, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM oracles WHERE name = ?", name)

	var oracle domain.Oracle
	if err := row.Scan(&oracle.Id, &oracle.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &oracle, nil
}

func (r *oracleRepo) FindByIds(ids []string) ([]domain.Oracle, error) {
	rows, err := r.connection.GetDB().Query("SELECT * FROM oracles WHERE id IN ?", ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var oracles []domain.Oracle
	for rows.Next() {
		var oracle domain.Oracle
		if err := rows.Scan(&oracle.Id, &oracle.Name); err != nil {
			return nil, err
		}
		oracles = append(oracles, oracle)
	}
	return oracles, nil
}

func (r *oracleRepo) FindByNames(names []string) ([]domain.Oracle, error) {
	rows, err := r.connection.GetDB().Query("SELECT * FROM oracles WHERE name IN ?", names)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var oracles []domain.Oracle
	for rows.Next() {
		var oracle domain.Oracle
		if err := rows.Scan(&oracle.Id, &oracle.Name); err != nil {
			return nil, err
		}
		oracles = append(oracles, oracle)
	}
	return oracles, nil
}

func (r *oracleRepo) Delete(id string) error {
	res, err := r.connection.GetDB().Exec("DELETE FROM oracles WHERE id = ?", id)
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
