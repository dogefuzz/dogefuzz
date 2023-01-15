package topic

import "github.com/dogefuzz/dogefuzz/pkg/interfaces"

type env interface {
	EventBus() interfaces.EventBus
}
