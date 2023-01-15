package topic

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

const INSTRUMENT_EXECUTION_TOPIC = "instrument:execution"

type instrumentExecutionTopic struct {
	eventBus interfaces.EventBus
}

func NewInstrumentExecutionTopic(e env) *instrumentExecutionTopic {
	return &instrumentExecutionTopic{eventBus: e.EventBus()}
}

func (t *instrumentExecutionTopic) Publish(e bus.InstrumentExecutionEvent) {
	t.eventBus.Publish(INSTRUMENT_EXECUTION_TOPIC, e)
}

func (t *instrumentExecutionTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(INSTRUMENT_EXECUTION_TOPIC, fn)
}

func (t *instrumentExecutionTopic) Unsubscribe(fn interface{}) {
	t.eventBus.Unsubscribe(INSTRUMENT_EXECUTION_TOPIC, fn)
}
