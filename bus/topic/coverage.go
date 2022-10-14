package topic

import (
	"github.com/gongbell/contractfuzzer/bus"
	"github.com/gongbell/contractfuzzer/bus/event"
)

const TASK_COVERAGE_TOPIC = "instrument:coverage"

type CoverageTopic interface {
	Publish(e event.CoverageEvent)
	Subscribe(fn interface{})
}

type DefaultCoverageTopic struct {
	eventBus bus.EventBus
}

func (t DefaultCoverageTopic) Init(eventBus bus.EventBus) DefaultCoverageTopic {
	t.eventBus = eventBus

	return t
}

func (t DefaultCoverageTopic) Publish(e event.TaskFinishEvent) {
	t.eventBus.Publish(TASK_COVERAGE_TOPIC, e)
}

func (t DefaultCoverageTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(TASK_COVERAGE_TOPIC, fn)
}
