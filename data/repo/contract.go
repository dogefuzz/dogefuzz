package repo

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContractRepo interface {
	Create(tx *gorm.DB, contract *entities.Contract) error
	Update(tx *gorm.DB, contract *entities.Contract) error
	Find(tx *gorm.DB, id string) (*entities.Contract, error)
	FindByTaskId(tx *gorm.DB, taskId string) (*entities.Contract, error)
}

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
	if err := tx.First(&contract, updatedContract.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotExists
		}
	}
	return tx.Model(&contract).Updates(updatedContract).Error
}

func (r *contractRepo) Find(tx *gorm.DB, id string) (*entities.Contract, error) {
	var contract entities.Contract
	if err := tx.First(&contract, id).Error; err != nil {
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
