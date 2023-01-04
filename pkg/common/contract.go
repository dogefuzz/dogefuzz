package common

type Contract struct {
	Name          string
	AbiDefinition string
	CompiledCode  string
	Address       string
}

func NewContract(name, abiDefinition, compiledCode string) *Contract {
	return &Contract{
		Name:          name,
		AbiDefinition: abiDefinition,
		CompiledCode:  compiledCode,
	}
}
