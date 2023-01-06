package entities

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type Transaction struct {
	Id                   string    `gorm:"primaryKey"`
	Timestamp            time.Time `gorm:"not null"`
	BlockchainHash       string    `gorm:"index;unique;not null"`
	TaskId               string    `gorm:"not null"`
	FunctionId           string    `gorm:"not null"`
	Inputs               string    `gorm:"not null"`
	DetectedWeaknesses   string
	ExecutedInstructions string
	DeltaCoverage        int64
	DeltaMinDistance     int64
	Status               common.TransactionStatus `gorm:"not null"`
}
