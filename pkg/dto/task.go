package dto

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type StartTaskRequest struct {
	ContractSource string              `json:"contractSource"`
	ContractName   string              `json:"contractName"`
	Arguments      []string            `json:"arguments"`
	Duration       string              `json:"duration"`
	Detectors      []common.OracleType `json:"detectors"`
	FuzzingType    common.FuzzingType  `json:"fuzzingType"`
}

type StartTaskResponse struct {
	TaskId string `json:"taskId"`
}

type NewTaskDTO struct {
	Arguments                      []string            `json:"arguments"`
	Duration                       time.Duration       `json:"duration"`
	StartTime                      time.Time           `json:"startTime"`
	DeploymentTime                 time.Time           `json:"deploymentTime"`
	Expiration                     time.Time           `json:"expiration"`
	Detectors                      []common.OracleType `json:"detectors"`
	FuzzingType                    common.FuzzingType  `json:"fuzzingType"`
	AggregatedExecutedInstructions []string            `json:"aggregatedExecutedInstructions"`
	Status                         common.TaskStatus   `json:"status"`
}

type TaskDTO struct {
	Id                             string              `json:"id"`
	Arguments                      []string            `json:"arguments"`
	Duration                       time.Duration       `json:"duration"`
	StartTime                      time.Time           `json:"startTime"`
	DeploymentTime                 time.Time           `json:"deploymentTime"`
	Expiration                     time.Time           `json:"expiration"`
	Detectors                      []common.OracleType `json:"detectors"`
	FuzzingType                    common.FuzzingType  `json:"fuzzingType"`
	AggregatedExecutedInstructions []string            `json:"aggregatedExecutedInstructions"`
	Status                         common.TaskStatus   `json:"status"`
}
