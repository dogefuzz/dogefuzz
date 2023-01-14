package topic

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
)

const TASK_INPUT_REQUEST_TOPIC = "task:input_request"

type taskInputRequestTopic struct {
	eventBus bus.EventBus
}

func NewTaskInputRequestTopic(e env) *taskInputRequestTopic {
	return &taskInputRequestTopic{eventBus: e.EventBus()}
}

func (t *taskInputRequestTopic) Publish(e bus.TaskInputRequestEvent) {
	t.eventBus.Publish(TASK_INPUT_REQUEST_TOPIC, e)
}

func (t *taskInputRequestTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(TASK_INPUT_REQUEST_TOPIC, fn)
}

func (t *taskInputRequestTopic) Unsubscribe(fn interface{}) {
	t.eventBus.Unsubscribe(TASK_INPUT_REQUEST_TOPIC, fn)
}
