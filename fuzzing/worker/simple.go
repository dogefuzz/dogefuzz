package worker

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"go.uber.org/zap"
)

type simpleFuzzingWorker struct {
	cancelChannel chan bool
	logger        *zap.Logger
	abiDir        string
	outDir        string
}

func NewSimpleFuzzingWorker(
	cancelChannel chan bool,
	logger *zap.Logger,
	abiDir string,
	outDir string,
) *simpleFuzzingWorker {
	return &simpleFuzzingWorker{
		cancelChannel: cancelChannel,
		logger:        logger,
		abiDir:        abiDir,
		outDir:        outDir,
	}
}

func (w *simpleFuzzingWorker) Start(taskId string, contracts []string, duration time.Duration) {
	w.logger.Info(fmt.Sprintf("Running a simple fuzzing for %-8v", duration))
	// go fuzz.Start(w.abiDir, w.outDir, taskId)
	<-w.cancelChannel
	// fuzz.G_stop <- true
}

func (w *simpleFuzzingWorker) GenerateInput(taskId string, contract *abi.ABI) {

}
