package domain

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type Task struct {
	Id         string
	Contract   string
	Expiration time.Time
	Detectors  string
	Status     common.TaskStatus
}
