package topic

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
)

const TASK_START_TOPIC = "task:start"

type taskStartTopic struct {
	eventBus bus.EventBus
}

func NewTaskStartTopic(e env) *taskStartTopic {
	return &taskStartTopic{eventBus: e.EventBus()}
}

func (t *taskStartTopic) Publish(e bus.TaskStartEvent) {
	t.eventBus.Publish(TASK_START_TOPIC, e)
}

func (t *taskStartTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(TASK_START_TOPIC, fn)
}

func (t *taskStartTopic) Unsubscribe(fn interface{}) {
	t.eventBus.Unsubscribe(TASK_START_TOPIC, fn)
}
