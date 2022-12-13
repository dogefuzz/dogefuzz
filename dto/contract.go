package dto

import "github.com/dogefuzz/dogefuzz/pkg/common"

type NewContractDTO struct {
	Source        string `json:"source"`
	CompiledCode  string `json:"compiledCode"`
	AbiDefinition string `json:"abiDefinition"`
	Name          string `json:"name"`
}

type ContractDTO struct {
	Id            string             `json:"id"`
	Address       string             `json:"address"`
	Source        string             `json:"source"`
	CompiledCode  string             `json:"compiledCode"`
	AbiDefinition string             `json:""`
	Name          string             `json:"name"`
	CFG           common.CFG         `json:"cfg"`
	DistanceMap   common.DistanceMap `json:"distanceMap"`
}
