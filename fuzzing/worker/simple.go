package worker

import (
	"fmt"
	"time"

	"github.com/dogefuzz/dogefuzz/fuzz"
	"go.uber.org/zap"
)

type SimpleFuzzingWorker struct {
	cancelChannel chan bool
	logger        *zap.Logger
	abiDir        string
	outDir        string
}

func (w SimpleFuzzingWorker) Init(
	cancelChannel chan bool,
	logger *zap.Logger,
	abiDir string,
	outDir string,
) SimpleFuzzingWorker {
	w.cancelChannel = cancelChannel
	w.logger = logger
	w.abiDir = abiDir
	w.outDir = outDir

	return w
}

func (w SimpleFuzzingWorker) Start(taskId string, contracts []string, duration time.Duration) {
	w.logger.Info(fmt.Sprintf("Running a simple fuzzing for %-8v", duration))
	go fuzz.Start(w.abiDir, w.outDir, taskId)
	<-w.cancelChannel
	fuzz.G_stop <- true
}
