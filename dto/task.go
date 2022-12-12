package dto

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type StartTaskRequest struct {
	Contract  string   `json:"contract"`
	Duration  string   `json:"duration"`
	Detectors []string `json:"detectors"`
}

type StartTaskResponse struct {
	TaskId string `json:"taskId"`
}

type NewTaskDTO struct {
	Contract   string            `json:"contract"`
	Expiration time.Time         `json:"expiration"`
	Detectors  []string          `json:"detectors"`
	Status     common.TaskStatus `json:"status"`
}

type TaskDTO struct {
	Id         string            `json:"id"`
	Contract   string            `json:"contract"`
	Expiration time.Time         `json:"expiration"`
	Detectors  []string          `json:"detectors"`
	Status     common.TaskStatus `json:"status"`
}
