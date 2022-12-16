package domain

import "github.com/dogefuzz/dogefuzz/pkg/common"

type Transaction struct {
	Id                   string
	BlockchainHash       string
	TaskId               string
	FunctionId           string
	Inputs               []string
	DetectedWeaknesses   string
	ExecutedInstructions string
	Status               common.TransactionStatus
}
