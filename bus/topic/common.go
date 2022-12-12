package topic

import "github.com/dogefuzz/dogefuzz/bus"

type env interface {
	EventBus() bus.EventBus
}
