package dto

type NewFunctionDTO struct {
	Name          string `json:"name"`
	NumberOfArgs  int64  `json:"numberOfArgs"`
	Payable       bool   `json:"payable"`
	IsConstructor bool   `json:"isConstructor"`
	ContractId    string `json:"contractId"`
}

type FunctionDTO struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	NumberOfArgs  int64  `json:"numberOfArgs"`
	Payable       bool   `json:"payable"`
	IsConstructor bool   `json:"isConstrutor"`
	ContractId    string `json:"contractId"`
}
