package entities

import "github.com/dogefuzz/dogefuzz/pkg/common"

type Function struct {
	Id           string            `gorm:"primaryKey"`
	Name         string            `gorm:"index;not null"`
	NumberOfArgs int64             `gorm:"not null"`
	Callable     bool              `gorm:"not null"`
	Type         common.MethodType `gorm:"not null"`
	ContractId   string            `gorm:"not null"`
	Transactions []Transaction
}
