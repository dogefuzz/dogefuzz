package entities

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type Task struct {
	Id          string
	ContractId  string
	Arguments   string
	StartTime   time.Time
	Expiration  time.Time
	Detectors   string
	FuzzingType common.FuzzingType
	Status      common.TaskStatus
}
