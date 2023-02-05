package fuzz

import (
	"github.com/dogefuzz/dogefuzz/pkg/solidity"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type blackboxFuzzer struct {
}

func NewBlackboxFuzzer() *blackboxFuzzer {
	return &blackboxFuzzer{}
}

func (f *blackboxFuzzer) GenerateInput(method abi.Method) ([]interface{}, error) {
	args := make([]interface{}, len(method.Inputs))
	for idx, input := range method.Inputs {
		handler, err := solidity.GetTypeHandler(input.Type)
		if err != nil {
			return nil, err
		}
		handler.Generate()
		args[idx] = handler.GetValue()
	}
	return args, nil
}
