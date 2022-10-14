package bus

import (
	evbus "github.com/asaskevich/EventBus"
)

type MemoryEventBus struct {
	InternalBus evbus.Bus
}

func (b MemoryEventBus) Init() (EventBus, error) {
	b.InternalBus = evbus.New()
	return b, nil
}

func (b MemoryEventBus) Subscribe(topic string, fn interface{}) {
	b.InternalBus.SubscribeAsync(topic, fn, true)
}

func (b MemoryEventBus) SubscribeOnce(topic string, fn interface{}) {
	b.InternalBus.SubscribeOnceAsync(topic, fn)
}

func (b MemoryEventBus) Publish(topic string, args ...interface{}) {
	b.InternalBus.Publish(topic, args...)
}
