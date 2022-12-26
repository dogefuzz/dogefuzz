package topic

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
)

const TASK_FINISH_TOPIC = "task:finish"

type taskFinishTopic struct {
	eventBus bus.EventBus
}

func NewTaskFinishTopic(e env) *taskFinishTopic {
	return &taskFinishTopic{eventBus: e.EventBus()}
}

func (t *taskFinishTopic) Publish(e bus.TaskFinishEvent) {
	t.eventBus.Publish(TASK_FINISH_TOPIC, e)
}

func (t *taskFinishTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(TASK_FINISH_TOPIC, fn)
}
