package repository

import (
	"database/sql"
	"errors"

	"github.com/gongbell/contractfuzzer/db"
	"github.com/gongbell/contractfuzzer/db/domain"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type ContractRepository interface {
	Create(contract *domain.Contract) error
	Find(id string) (*domain.Contract, error)
	FindByName(name string) (*domain.Contract, error)
	FindByAddress(address string) (*domain.Contract, error)
	Delete(id string) error
}

type ContractSQLiteRepository struct {
	manager db.Manager
}

func (r ContractSQLiteRepository) Init(manager db.Manager) (ContractSQLiteRepository, error) {
	r.manager = manager
	return r, nil
}

func (r ContractSQLiteRepository) Create(contract *domain.Contract) error {
	contract.Id = uuid.NewString()
	_, err := r.manager.GetDB().Exec("INSERT INTO contracts(id, name, address) values(?, ?, ?)", contract.Id, contract.Name, contract.Address)
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

func (r ContractSQLiteRepository) Find(id string) (*domain.Contract, error) {
	row := r.manager.GetDB().QueryRow("SELECT * FROM contracts WHERE id = ?", id)

	var contract domain.Contract
	if err := row.Scan(&contract.Id, &contract.Name, &contract.Address); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &contract, nil
}

func (r ContractSQLiteRepository) FindByName(name string) (*domain.Contract, error) {
	row := r.manager.GetDB().QueryRow("SELECT * FROM contracts WHERE name = ?", name)

	var contract domain.Contract
	if err := row.Scan(&contract.Id, &contract.Name, &contract.Address); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &contract, nil
}

func (r ContractSQLiteRepository) FindByAddress(address string) (*domain.Contract, error) {
	row := r.manager.GetDB().QueryRow("SELECT * FROM contracts WHERE address = ?", address)

	var contract domain.Contract
	if err := row.Scan(&contract.Id, &contract.Name, &contract.Address); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &contract, nil
}

func (r ContractSQLiteRepository) Delete(id string) error {
	res, err := r.manager.GetDB().Exec("DELETE FROM contracts WHERE id = ?", id)
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
