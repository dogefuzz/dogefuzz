package entities

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type Transaction struct {
	Id                   string
	Timestamp            time.Time
	BlockchainHash       string
	TaskId               string
	FunctionId           string
	Inputs               []string
	DetectedWeaknesses   string
	ExecutedInstructions string
	DeltaCoverage        int64
	DeltaMinDistance     int64
	Status               common.TransactionStatus
}
