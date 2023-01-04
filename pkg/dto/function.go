package dto

type NewFunctionDTO struct {
	Name         string
	NumberOfArgs int64
	Payable      bool
}

type FunctionDTO struct {
	Id           string
	Name         string
	NumberOfArgs int64
	Payable      bool
}
