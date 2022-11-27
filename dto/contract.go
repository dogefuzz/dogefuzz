package dto

type NewContractDTO struct {
	Source        string `json:"source"`
	CompiledCode  string `json:"compiledCode"`
	AbiDefinition string `json:"abiDefinition"`
	Name          string `json:"name"`
}

type ContractDTO struct {
	Id            string `json:"id"`
	Address       string `json:"address"`
	Source        string `json:"source"`
	CompiledCode  string `json:"compiledCode"`
	AbiDefinition string `json:""`
	Name          string `json:"name"`
}
