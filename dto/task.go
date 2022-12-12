package dto

import "time"

type StartTaskRequest struct {
	Contract  string   `json:"contract"`
	Duration  string   `json:"duration"`
	Detectors []string `json:"detectors"`
}

type StartTaskResponse struct {
	TaskId string `json:"taskId"`
}

type NewTaskDTO struct {
	Contract  string        `json:"contract"`
	Duration  time.Duration `json:"duration"`
	Detectors []string      `json:"detectors"`
}

type TaskDTO struct {
	Id        string        `json:"id"`
	Contract  string        `json:"contract"`
	Duration  time.Duration `json:"duration"`
	Detectors []string      `json:"detectors"`
}


