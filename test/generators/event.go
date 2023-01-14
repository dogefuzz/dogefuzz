package generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
)

func InstrumentExecutionEventGen() bus.InstrumentExecutionEvent {
	return bus.InstrumentExecutionEvent{TransactionId: gofakeit.LetterN(255)}
}

func TaskFinishEventGen() bus.TaskFinishEvent {
	return bus.TaskFinishEvent{TaskId: gofakeit.LetterN(255)}
}

func TaskInputRequestEventGen() bus.TaskInputRequestEvent {
	return bus.TaskInputRequestEvent{TaskId: gofakeit.LetterN(255)}
}

func TaskStartEventGen() bus.TaskStartEvent {
	return bus.TaskStartEvent{TaskId: gofakeit.LetterN(255)}
}
