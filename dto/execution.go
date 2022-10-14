package dto

type NewExecutionDTO struct {
	Name         string   `json:"name"`
	Input        string   `json:"input"`
	Instructions []uint64 `json:"instructions"`
	TxHash       string   `json:"txHash"`
}
