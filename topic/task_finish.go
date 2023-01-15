package topic

import (
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

const TASK_FINISH_TOPIC = "task:finish"

type taskFinishTopic struct {
	eventBus interfaces.EventBus
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

func (t *taskFinishTopic) Unsubscribe(fn interface{}) {
	t.eventBus.Unsubscribe(TASK_FINISH_TOPIC, fn)
}
