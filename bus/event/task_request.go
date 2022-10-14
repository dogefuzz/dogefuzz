package event

import "time"

type TaskRequestEvent struct {
	TaskId      string
	Duration    time.Duration
	Contracts   []string
	FuzzingType string
}
