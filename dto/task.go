package dto

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type StartTaskRequest struct {
	Contract  string              `json:"contract"`
	Duration  string              `json:"duration"`
	Detectors []common.OracleType `json:"detectors"`
}

type StartTaskResponse struct {
	TaskId string `json:"taskId"`
}

type NewTaskDTO struct {
	ContractId string              `json:"contractId"`
	Expiration time.Time           `json:"expiration"`
	Detectors  []common.OracleType `json:"detectors"`
	Status     common.TaskStatus   `json:"status"`
}

type TaskDTO struct {
	Id         string              `json:"id"`
	ContractId string              `json:"contractId"`
	Expiration time.Time           `json:"expiration"`
	Detectors  []common.OracleType `json:"detectors"`
	Status     common.TaskStatus   `json:"status"`
}
