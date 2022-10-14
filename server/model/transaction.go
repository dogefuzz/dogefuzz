package model

type TransactionCreateRequest struct {
	TaskId          string `json:"taskId"`
	BlockchainHash  string `json:"blockchainHash"`
	ContractAddress string `json:"contractAddress"`
}

type TransactionCreateResponse struct {
	TransactionId string `json:"transactionId"`
}
