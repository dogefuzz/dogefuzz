package fuzzing

import (
	"errors"
	"fmt"
	"time"

	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/event"
	"github.com/dogefuzz/dogefuzz/fuzzing/worker"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"go.uber.org/zap"
)

var (
	ErrWorkerTypeNotFound       = errors.New("worker type not found")
	ErrWorkerTypeNotImplemented = errors.New("worker type not implemented")
)

type FuzzingTask struct {
	TaskCancelChannel  chan bool
	TimerCancelChannel chan bool
}

type FuzzingMaster interface {
	StartFuzzer(taskId string, contracts []string, duration time.Duration)
	StartTimer(taskId string, duration time.Duration)
}

type fuzzingMaster struct {
	Logger         *zap.Logger
	EventBus       bus.EventBus
	AbiDir         string
	OutDir         string
	Workers        *common.ConcurrentMap
	RequestChannel chan string
}

func NewFuzzingMaster(
	logger *zap.Logger,
	eventBus bus.EventBus,
	abiDir string,
	outDir string,
) *fuzzingMaster {
	return &fuzzingMaster{
		Logger:   logger,
		EventBus: eventBus,
		AbiDir:   abiDir,
		OutDir:   outDir,
		Workers:  &common.ConcurrentMap{},
	}
}

func (m *fuzzingMaster) StartFuzzer(e event.TaskRequestEvent) {
	m.Logger.Info(fmt.Sprintf("Running fuzzing task %s for %-8v", e.TaskId, e.Duration))

	taskCancelChannel := make(chan bool)
	timerCancelChannel := make(chan bool)

	worker, err := m.getWorkerType(e.FuzzingType, taskCancelChannel)
	if err != nil {
		m.Logger.Error(fmt.Sprintf("Error while starting worker: %s", e.TaskId))
		return
	}

	workerInfo := FuzzingTask{}
	workerInfo.TaskCancelChannel = taskCancelChannel
	workerInfo.TimerCancelChannel = timerCancelChannel
	m.Workers.Add(e.TaskId, workerInfo)

	go m.StartTimer(e.TaskId, e.Duration, timerCancelChannel)
	worker.Start(e.TaskId, e.Contracts, e.Duration)
}

func (m *fuzzingMaster) StartTimer(taskId string, duration time.Duration, timerCancelChannel chan bool) {
	m.Logger.Info(fmt.Sprintf("Start timer for %-8v", duration))

	timer := time.NewTimer(duration)
loop:
	for {
		select {
		case <-timer.C:
			m.Logger.Info(fmt.Sprintf("Stopping fuzzer task %s", taskId))
			m.EventBus.Publish("task:finish", taskId)
			m.Logger.Info(fmt.Sprintf("Stopping fuzzing timer %s", taskId))
			break loop
		case <-timerCancelChannel:
			m.Logger.Info(fmt.Sprintf("Stopping fuzzing timer %s", taskId))
			break loop
		}
	}
}

func (m *fuzzingMaster) getWorkerType(fuzzingType string, taskCancelChannel chan bool) (worker.FuzzingWorker, error) {
	var fuzzingWorker worker.FuzzingWorker
	switch fuzzingType {
	case BLACKBOX_FUZZING:
		fuzzingWorker = worker.NewSimpleFuzzingWorker(taskCancelChannel, m.Logger, m.AbiDir, m.OutDir)
	case GREYBOX_FUZZING:
		return nil, ErrWorkerTypeNotImplemented
	case DIRECTED_GREYBOX_FUZZING:
		return nil, ErrWorkerTypeNotImplemented
	default:
		return nil, ErrWorkerTypeNotFound
	}
	return fuzzingWorker, nil
}
