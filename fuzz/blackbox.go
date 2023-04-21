package fuzz

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type blackboxFuzzer struct {
	solidityService interfaces.SolidityService
	functionService interfaces.FunctionService
	contractService interfaces.ContractService
}

func NewBlackboxFuzzer(e env) *blackboxFuzzer {
	return &blackboxFuzzer{
		solidityService: e.SolidityService(),
		functionService: e.FunctionService(),
		contractService: e.ContractService(),
	}
}

func (f *blackboxFuzzer) GenerateInput(functionId string) ([]interface{}, error) {
	function, err := f.functionService.Get(functionId)
	if err != nil {
		return nil, err
	}

	contract, err := f.contractService.Get(function.ContractId)
	if err != nil {
		return nil, err
	}

	abiDefinition, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		return nil, err
	}
	method := abiDefinition.Methods[function.Name]

	args := make([]interface{}, len(method.Inputs))
	for idx, input := range method.Inputs {
		handler, err := f.solidityService.GetTypeHandlerWithContext(input.Type)
		if err != nil {
			return nil, err
		}
		handler.Generate()
		args[idx] = handler.GetValue()
	}
	return args, nil
}
