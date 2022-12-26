package bus

import (
	evbus "github.com/asaskevich/EventBus"
)

type memoryEventBus struct {
	internalBus evbus.Bus
}

func NewMemoryEventBus() *memoryEventBus {
	return &memoryEventBus{
		internalBus: evbus.New(),
	}
}

func (b *memoryEventBus) Subscribe(topic string, fn interface{}) {
	b.internalBus.SubscribeAsync(topic, fn, true)
}

func (b *memoryEventBus) SubscribeOnce(topic string, fn interface{}) {
	b.internalBus.SubscribeOnceAsync(topic, fn)
}

func (b *memoryEventBus) Publish(topic string, args ...interface{}) {
	b.internalBus.Publish(topic, args...)
}
