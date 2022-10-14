package model

type FuzzStartRequest struct {
	Contracts []string `json:"contracts"`
	Duration  string   `json:"duration"`
	Detectors []string `json:"detectors"`
}

type FuzzStartResponse struct {
	TaskId string `json:"taskId"`
}
