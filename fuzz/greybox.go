package fuzz

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/pkg/solidity"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type greyboxFuzzer struct {
	powerSchedule interfaces.PowerSchedule
}

func NewGreyboxFuzzer(e env) *greyboxFuzzer {
	return &greyboxFuzzer{
		powerSchedule: e.PowerSchedule(),
	}
}

func (f *greyboxFuzzer) GenerateInput(method abi.Method) ([]interface{}, error) {
	seedsList, err := f.powerSchedule.RequestSeeds(method, common.COVERAGE_BASED_STRATEGY)
	if err != nil {
		return nil, err
	}

	chosenSeeds := common.RandomChoice(seedsList)

	inputs := make([]interface{}, len(method.Inputs))
	for inputsIdx, inputDefinition := range method.Inputs {
		handler, err := solidity.GetTypeHandler(inputDefinition.Type)
		if err != nil {
			return nil, err
		}
		handler.SetValue(chosenSeeds[inputsIdx])
		mutationFunction := common.RandomChoice(handler.GetMutators())
		mutationFunction()
		inputs[inputsIdx] = handler.GetValue()
	}
	return inputs, nil
}
