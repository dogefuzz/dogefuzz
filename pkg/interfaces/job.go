package interfaces

import "context"

type CronJob interface {
	Id() string
	CronConfig() string
	Handler()
}

type Scheduler interface {
	Start()
	Shutdown() context.Context
}
