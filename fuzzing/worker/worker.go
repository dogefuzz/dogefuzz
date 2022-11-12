package worker

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type FuzzingWorker interface {
	Start(taskId string, contracts []string, duration time.Duration)
	GenerateInput(taskId string, contract *abi.ABI)
}
