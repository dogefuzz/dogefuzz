package topic

import (
	"github.com/dogefuzz/dogefuzz/bus"
)

const TASK_FINISH_TOPIC = "task:finish"

type taskFinishTopic struct {
	eventBus bus.EventBus
}

func NewTaskFinishTopic(eventBus bus.EventBus) *taskFinishTopic {
	return &taskFinishTopic{eventBus: eventBus}
}

func (t *taskFinishTopic) Publish(e bus.TaskFinishEvent) {
	t.eventBus.Publish(TASK_FINISH_TOPIC, e)
}

func (t *taskFinishTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(TASK_FINISH_TOPIC, fn)
}
