package topic

import (
	"github.com/dogefuzz/dogefuzz/bus"
)

const TASK_REQUEST_TOPIC = "task:request"

type taskStartTopic struct {
	eventBus bus.EventBus
}

func NewTaskStartTopic(e env) *taskStartTopic {
	return &taskStartTopic{eventBus: e.EventBus()}
}

func (t *taskStartTopic) Publish(e bus.TaskStartEvent) {
	t.eventBus.Publish(TASK_REQUEST_TOPIC, e)
}

func (t *taskStartTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(TASK_REQUEST_TOPIC, fn)
}
