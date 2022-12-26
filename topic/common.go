package topic

import "github.com/dogefuzz/dogefuzz/pkg/bus"

type env interface {
	EventBus() bus.EventBus
}
