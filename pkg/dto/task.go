package dto

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type StartTaskRequest struct {
	Contract  string              `json:"contract"`
	Arguments []string            `json:"arguments"`
	Duration  string              `json:"duration"`
	Detectors []common.OracleType `json:"detectors"`
}

type StartTaskResponse struct {
	TaskId string `json:"taskId"`
}

type NewTaskDTO struct {
	ContractId  string              `json:"contractId"`
	Arguments   []string            `json:"arguments"`
	Expiration  time.Time           `json:"expiration"`
	Detectors   []common.OracleType `json:"detectors"`
	FuzzingType common.FuzzingType  `json:"fuzzingType"`
	Status      common.TaskStatus   `json:"status"`
}

type TaskDTO struct {
	Id                             string              `json:"id"`
	ContractId                     string              `json:"contractId"`
	Arguments                      []string            `json:"arguments"`
	StartTime                      time.Time           `json:"startTime"`
	Expiration                     time.Time           `json:"expiration"`
	Detectors                      []common.OracleType `json:"detectors"`
	FuzzingType                    common.FuzzingType  `json:"fuzzingType"`
	AggregatedExecutedInstructions []string            `json:"aggregatedExecutedInstructions"`
	Status                         common.TaskStatus   `json:"status"`
}
