package dto

import "github.com/dogefuzz/dogefuzz/pkg/common"

type NewContractDTO struct {
	TaskId             string                `json:"taskId"`
	Status             common.ContractStatus `json:"status"`
	Source             string                `json:"source"`
	DeploymentBytecode string                `json:"deploymentBytecode"`
	RuntimeBytecode    string                `json:"runtimeBytecode"`
	AbiDefinition      string                `json:"abiDefinition"`
	Name               string                `json:"name"`
}

type NewContractWithIdDTO struct {
	Id                 string                `json:"id"`
	TaskId             string                `json:"taskId"`
	Status             common.ContractStatus `json:"status"`
	Source             string                `json:"source"`
	DeploymentBytecode string                `json:"deploymentBytecode"`
	RuntimeBytecode    string                `json:"runtimeBytecode"`
	AbiDefinition      string                `json:"abiDefinition"`
	Name               string                `json:"name"`
}

type ContractDTO struct {
	Id                 string                `json:"id"`
	TaskId             string                `json:"taskId"`
	Status             common.ContractStatus `json:"status"`
	Address            string                `json:"address"`
	Source             string                `json:"source"`
	DeploymentBytecode string                `json:"deploymentBytecode"`
	RuntimeBytecode    string                `json:"runtimeBytecode"`
	AbiDefinition      string                `json:"abiDefinition"`
	Name               string                `json:"name"`
	CFG                common.CFG            `json:"cfg"`
	DistanceMap        common.DistanceMap    `json:"distanceMap"`
}
