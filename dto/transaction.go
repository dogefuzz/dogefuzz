package dto

type NewTransactionDTO struct {
	TaskId          string `json:"taskId"`
	BlockchainHash  string `json:"blockchainHash"`
	ContractAddress string `json:"contractAddress"`
}

type TransactionDTO struct {
	TransactionId string `json:"transactionId"`
}
