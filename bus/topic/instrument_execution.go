package topic

import (
	"github.com/dogefuzz/dogefuzz/bus"
)

const INSTRUMENT_EXECUTION_TOPIC = "instrument:execution"

type instrumentExecutionTopic struct {
	eventBus bus.EventBus
}

func NewInstrumentExecutionTopic(eventBus bus.EventBus) *instrumentExecutionTopic {
	return &instrumentExecutionTopic{eventBus: eventBus}
}

func (t *instrumentExecutionTopic) Publish(e bus.InstrumentExecutionEvent) {
	t.eventBus.Publish(INSTRUMENT_EXECUTION_TOPIC, e)
}

func (t *instrumentExecutionTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(INSTRUMENT_EXECUTION_TOPIC, fn)
}
