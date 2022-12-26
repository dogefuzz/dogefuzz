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

func (f *blackboxFuzzer) GenerateInput(method abi.Method) []interface{} {
	args := make([]interface{}, len(method.Inputs))
	for _, input := range method.Inputs {
		handler, _ := solidity.GetTypeHandler(input.Type)
		handler.Generate()
		args = append(args, handler.GetValue())
	}
	return args
}
