package entities

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type Task struct {
	Id                             string             `gorm:"primaryKey"`
	Arguments                      string             `gorm:"not null"`
	StartTime                      time.Time          `gorm:"not null"`
	Expiration                     time.Time          `gorm:"not null"`
	Detectors                      string             `gorm:"not null"`
	FuzzingType                    common.FuzzingType `gorm:"not null"`
	AggregatedExecutedInstructions string             `gorm:"not null"`
	Status                         common.TaskStatus  `gorm:"not null"`
	Contract                       Contract
	Transactions                   []Transaction
}
