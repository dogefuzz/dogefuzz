package common

type Contract struct {
	Name               string
	AbiDefinition      string
	DeploymentBytecode string
	RuntimeBytecode    string
	Address            string
}

func NewContract(name, abiDefinition, deploymentBytecode string, runtimeBytecode string) *Contract {
	return &Contract{
		Name:               name,
		AbiDefinition:      abiDefinition,
		DeploymentBytecode: deploymentBytecode,
		RuntimeBytecode:    runtimeBytecode,
	}
}
