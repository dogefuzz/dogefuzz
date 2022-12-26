package fuzz

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/solidity"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type directedGreyboxFuzzer struct {
	powerSchedule PowerSchedule
}

func NewDirectedGreyboxFuzzer(e env) *directedGreyboxFuzzer {
	return &directedGreyboxFuzzer{
		powerSchedule: e.PowerSchedule(),
	}
}

func (f *directedGreyboxFuzzer) GenerateInput(method abi.Method) ([]interface{}, error) {
	seedsList, err := f.powerSchedule.RequestSeeds(method, common.DISTANCE_BASED_STRATEGY)
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
