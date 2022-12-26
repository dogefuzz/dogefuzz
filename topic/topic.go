package topic

type Topic[E any] interface {
	Publish(e E)
	Subscribe(fn interface{})
}
