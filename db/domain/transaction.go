package domain

import "strings"

type Transaction struct {
	Id                 string
	BlockchainHash     string
	TaskId             string
	ContractId         string
	DetectedWeaknesses string
}

func (t Transaction) GetDetectedWeaknesses() []string {
	return strings.Split(t.DetectedWeaknesses, ";")
}

func (t Transaction) SetDetectedWeaknesses(weaknesses []string) {
	t.DetectedWeaknesses = strings.Join(weaknesses, ";")
}
