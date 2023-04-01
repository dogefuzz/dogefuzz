package entities

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type Task struct {
	Id                             string `gorm:"primaryKey"`
	Arguments                      string `gorm:"not null"`
	Duration                       time.Duration
	StartTime                      time.Time `gorm:"not null"`
	DeploymentTime                 time.Time
	Expiration                     time.Time          `gorm:"not null"`
	Detectors                      string             `gorm:"not null"`
	FuzzingType                    common.FuzzingType `gorm:"not null"`
	AggregatedExecutedInstructions string
	Status                         common.TaskStatus `gorm:"not null"`
	Contract                       Contract
	Transactions                   []Transaction
}
