package fuzz

import (
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type blackboxFuzzer struct {
	solidityService interfaces.SolidityService
}

func NewBlackboxFuzzer(e env) *blackboxFuzzer {
	return &blackboxFuzzer{
		solidityService: e.SolidityService(),
	}
}

func (f *blackboxFuzzer) GenerateInput(method abi.Method) ([]interface{}, error) {
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
