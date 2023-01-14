package repo

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type functionRepo struct {
}

func NewFunctionRepo(e Env) *functionRepo {
	return &functionRepo{}
}

func (r *functionRepo) Get(tx *gorm.DB, id string) (*entities.Function, error) {
	var function entities.Function
	if err := tx.First(&function, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &function, nil
}

func (r *functionRepo) Create(tx *gorm.DB, function *entities.Function) error {
	function.Id = uuid.NewString()
	return tx.Create(function).Error
}

func (r *functionRepo) FindByContractId(tx *gorm.DB, contractId string) ([]entities.Function, error) {
	var functions []entities.Function
	if err := tx.Where("contract_id = ?", contractId).Find(&functions).Error; err != nil {
		return nil, err
	}
	return functions, nil
}

func (r *functionRepo) FindConstructorByContractId(tx *gorm.DB, contractId string) (*entities.Function, error) {
	var function entities.Function
	if err := tx.Where("is_constructor = ?", true).Where("contract_id = ?", contractId).First(&function).Error; err != nil {
		return nil, err
	}
	return &function, nil
}
