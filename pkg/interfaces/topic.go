package interfaces

type Topic[E any] interface {
	Publish(e E)
	Subscribe(fn interface{})
	Unsubscribe(fn interface{})
}
