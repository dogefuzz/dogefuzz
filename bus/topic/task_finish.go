package topic

import (
	"github.com/gongbell/contractfuzzer/bus"
	"github.com/gongbell/contractfuzzer/bus/event"
)

const TASK_FINISH_TOPIC = "task:finish"

type TaskFinishTopic interface {
	Publish(e event.TaskFinishEvent)
	Subscribe(fn interface{})
}

type DefaultTaskFinishTopic struct {
	eventBus bus.EventBus
}

func (t DefaultTaskFinishTopic) Init(eventBus bus.EventBus) DefaultTaskFinishTopic {
	t.eventBus = eventBus

	return t
}

func (t DefaultTaskFinishTopic) Publish(e event.TaskFinishEvent) {
	t.eventBus.Publish(TASK_FINISH_TOPIC, e)
}

func (t DefaultTaskFinishTopic) Subscribe(fn interface{}) {
	t.eventBus.Subscribe(TASK_FINISH_TOPIC, fn)
}
