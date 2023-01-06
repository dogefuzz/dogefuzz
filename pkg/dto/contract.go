package dto

import "github.com/dogefuzz/dogefuzz/pkg/common"

type NewContractDTO struct {
	TaskId        string `gorm:"foreignKey;not null"`
	Source        string `json:"source"`
	CompiledCode  string `json:"compiledCode"`
	AbiDefinition string `json:"abiDefinition"`
	Name          string `json:"name"`
}

type ContractDTO struct {
	Id            string             `json:"id"`
	TaskId        string             `gorm:"foreignKey;not null"`
	Address       string             `json:"address"`
	Source        string             `json:"source"`
	CompiledCode  string             `json:"compiledCode"`
	AbiDefinition string             `json:"abiDefinition"`
	Name          string             `json:"name"`
	CFG           common.CFG         `json:"cfg"`
	DistanceMap   common.DistanceMap `json:"distanceMap"`
}