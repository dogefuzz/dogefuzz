package interfaces

type EventBus interface {
	Subscribe(topic string, fn interface{})
	SubscribeOnce(topic string, fn interface{})
	Unsubscribe(topic string, fn interface{})
	Publish(topic string, args ...interface{})
}
