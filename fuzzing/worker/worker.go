package worker

import "time"

type FuzzingWorker interface {
	Start(taskId string, contracts []string, duration time.Duration)
}
