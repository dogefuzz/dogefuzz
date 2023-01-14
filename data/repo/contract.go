package repo

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type contractRepo struct {
}

func NewContractRepo(e Env) *contractRepo {
	return &contractRepo{}
}

func (r *contractRepo) Create(tx *gorm.DB, contract *entities.Contract) error {
	contract.Id = uuid.NewString()
	return tx.Create(contract).Error
}

func (r *contractRepo) Update(tx *gorm.DB, updatedContract *entities.Contract) error {
	var contract entities.Contract
	if err := tx.First(&contract, "id = ?", updatedContract.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotExists
		}
		return err
	}
	return tx.Model(&contract).Updates(updatedContract).Error
}

func (r *contractRepo) Find(tx *gorm.DB, id string) (*entities.Contract, error) {
	var contract entities.Contract
	if err := tx.First(&contract, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &contract, nil
}

func (r *contractRepo) FindByTaskId(tx *gorm.DB, taskId string) (*entities.Contract, error) {
	var contract entities.Contract
	if err := tx.Where("task_id = ?", taskId).First(&contract).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &contract, nil
}
