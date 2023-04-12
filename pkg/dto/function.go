package dto

import "github.com/dogefuzz/dogefuzz/pkg/common"

type NewFunctionDTO struct {
	Name         string            `json:"name"`
	NumberOfArgs int64             `json:"numberOfArgs"`
	Callable     bool              `json:"callable"`
	Type         common.MethodType `json:"type"`
	ContractId   string            `json:"contractId"`
}

type FunctionDTO struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	NumberOfArgs int64             `json:"numberOfArgs"`
	Callable     bool              `json:"callable"`
	Type         common.MethodType `json:"type"`
	ContractId   string            `json:"contractId"`
}
