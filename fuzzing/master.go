package fuzzing

import (
	"errors"
	"fmt"
	"time"

	"github.com/gongbell/contractfuzzer/bus"
	"github.com/gongbell/contractfuzzer/bus/event"
	"github.com/gongbell/contractfuzzer/fuzzing/worker"
	"github.com/gongbell/contractfuzzer/pkg/common"
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

type DefaultFuzzingMaster struct {
	Logger   *zap.Logger
	EventBus bus.EventBus
	AbiDir   string
	OutDir   string
	Workers  *common.ConcurrentMap
	RequestChannel chan string
}

func (m DefaultFuzzingMaster) Init(
	logger *zap.Logger,
	eventBus bus.EventBus,
	abiDir string,
	outDir string,
) DefaultFuzzingMaster {
	m.Logger = logger
	m.EventBus = eventBus
	m.AbiDir = abiDir
	m.OutDir = outDir
	m.Workers = &common.ConcurrentMap{}

	return m
}



func (m DefaultFuzzingMaster) startFuzzer(e event.TaskRequestEvent) {
	m.Logger.Info(fmt.Sprintf("Running fuzzing task %s for %-8v", e.TaskId, e.Duration))

	taskCancelChannel := make(chan bool, 0)
	timerCancelChannel := make(chan bool, 0)

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

func (m DefaultFuzzingMaster) StartTimer(taskId string, duration time.Duration, timerCancelChannel chan bool) {
	m.Logger.Info(fmt.Sprintf("Start timer for %-8v", duration))

	timer := time.NewTimer(duration)
	for {
		select {
		case <-timer.C:
			m.Logger.Info(fmt.Sprintf("Stopping fuzzer task %s", taskId))
			m.EventBus.Publish("task:finish", taskId)
			m.Logger.Info(fmt.Sprintf("Stopping fuzzing timer %s", taskId))
			break
		case <-timerCancelChannel:
			m.Logger.Info(fmt.Sprintf("Stopping fuzzing timer %s", taskId))
			break
		}
	}
}

func (m DefaultFuzzingMaster) getWorkerType(fuzzingType string, taskCancelChannel chan bool) (worker.FuzzingWorker, error) {
	var fuzzingWorker worker.FuzzingWorker
	switch fuzzingType {
	case BLACKBOX_FUZZING:
		fuzzingWorker = new(worker.SimpleFuzzingWorker).Init(taskCancelChannel, m.Logger, m.AbiDir, m.OutDir)
	case GREYBOX_FUZZING:
		return nil, ErrWorkerTypeNotImplemented
	case DIRECTED_GREYBOX_FUZZING:
		return nil, ErrWorkerTypeNotImplemented
	default:
		return nil, ErrWorkerTypeNotFound
	}
	return fuzzingWorker, nil
}
