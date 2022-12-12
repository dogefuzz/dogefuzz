package domain

type Transaction struct {
	Id                   string
	BlockchainHash       string
	TaskId               string
	ContractId           string
	DetectedWeaknesses   string
	ExecutedInstructions string
}
