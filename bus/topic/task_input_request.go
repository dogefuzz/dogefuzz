package topic

import (
	"github.com/dogefuzz/dogefuzz/bus"
)

const TASK_INPUT_REQUEST_TOPIC = "task:input_request"

type taskInputRequestTopic struct {
	eventBus bus.EventBus
}

func NewTaskInputRequestTopic(eventBus bus.EventBus) *taskInputRequestTopic {
	return &taskInputRequestTopic{eventBus: eventBus}
}

func (t *taskInputRequestTopic) Publish(e bus.TaskInputRequestEvent) {
	t.eventBus.Publish(INSTRUMENT_EXECUTION_TOPIC, e)
}

func (t *taskInputRequestTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(INSTRUMENT_EXECUTION_TOPIC, fn)
}
