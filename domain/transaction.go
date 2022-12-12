package domain

import "github.com/dogefuzz/dogefuzz/pkg/common"

type Transaction struct {
	Id                   string
	BlockchainHash       string
	TaskId               string
	ContractId           string
	DetectedWeaknesses   string
	ExecutedInstructions string
	Status               common.TransactionStatus
}
