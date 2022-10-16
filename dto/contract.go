package dto

type NewContractDTO struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type ContractDTO struct {
	Id      string `json:"id"`
	Address string `json:"address"`
	Source  string `json:"source"`
	Name    string `json:"name"`
}
