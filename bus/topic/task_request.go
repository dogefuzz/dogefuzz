package topic

import (
	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/event"
)

const TASK_REQUEST_TOPIC = "task:request"

type TaskRequestTopic struct {
	eventBus bus.EventBus
}

func (t TaskRequestTopic) Init(eventBus bus.EventBus) {
	t.eventBus = eventBus
}

func (t TaskRequestTopic) Publish(e event.TaskRequestEvent) {
	t.eventBus.Publish(TASK_REQUEST_TOPIC, e)
}

func (t TaskRequestTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(TASK_REQUEST_TOPIC, fn)
}
