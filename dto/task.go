package dto

type NewTaskDTO struct {
	Contracts []string `json:"contracts"`
	Duration  string   `json:"duration"`
	Detectors []string `json:"detectors"`
}

type TaskDTO struct {
	NewTaskDTO
	TaskId string `json:"taskId"`
}
