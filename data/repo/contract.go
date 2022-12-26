package repo

import (
	"database/sql"
	"errors"

	"github.com/dogefuzz/dogefuzz/data"
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type ContractRepo interface {
	Create(contract *entities.Contract) error
	Find(id string) (*entities.Contract, error)
	FindByName(name string) (*entities.Contract, error)
	FindByAddress(address string) (*entities.Contract, error)
	Delete(id string) error
}

type contractRepo struct {
	connection data.Connection
}

func NewContractRepo(e Env) *contractRepo {
	return &contractRepo{connection: e.DbConnection()}
}

func (r *contractRepo) Create(contract *entities.Contract) error {
	contract.Id = uuid.NewString()
	_, err := r.connection.GetDB().Exec("INSERT INTO contracts(id, name, source, address) values(?, ?, ?)", contract.Id, contract.Name, contract.Source, contract.Address)
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

func (r *contractRepo) Find(id string) (*entities.Contract, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM contracts WHERE id = ?", id)

	var contract entities.Contract
	if err := row.Scan(&contract.Id, &contract.Name, &contract.Source, &contract.Address); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &contract, nil
}

func (r *contractRepo) FindByName(name string) (*entities.Contract, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM contracts WHERE name = ?", name)

	var contract entities.Contract
	if err := row.Scan(&contract.Id, &contract.Name, &contract.Source, &contract.Address); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &contract, nil
}

func (r *contractRepo) FindByAddress(address string) (*entities.Contract, error) {
	row := r.connection.GetDB().QueryRow("SELECT * FROM contracts WHERE address = ?", address)

	var contract entities.Contract
	if err := row.Scan(&contract.Id, &contract.Name, &contract.Source, &contract.Address); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &contract, nil
}

func (r *contractRepo) Delete(id string) error {
	res, err := r.connection.GetDB().Exec("DELETE FROM contracts WHERE id = ?", id)
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
