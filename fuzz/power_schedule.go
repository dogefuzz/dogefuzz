package fuzz

import (
	"errors"
	"strings"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var ErrInvalidStrategy = errors.New("the provided strategy is not valid")
var ErrSeedsListInvalid = errors.New("the provided seeds list is invalid")

type powerSchedule struct {
	cfg                *config.Config
	transactionService interfaces.TransactionService
	solidityService    interfaces.SolidityService
	functionService    interfaces.FunctionService
	contractService    interfaces.ContractService
}

func NewPowerSchedule(e env) *powerSchedule {
	return &powerSchedule{
		cfg:                e.Config(),
		transactionService: e.TransactionService(),
		solidityService:    e.SolidityService(),
		functionService:    e.FunctionService(),
		contractService:    e.ContractService(),
	}
}

func (s *powerSchedule) RequestSeeds(functionId string, strategy common.PowerScheduleStrategy) ([][]interface{}, error) {
	function, err := s.functionService.Get(functionId)
	if err != nil {
		return nil, err
	}

	contract, err := s.contractService.Get(function.ContractId)
	if err != nil {
		return nil, err
	}

	abiDefinition, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		return nil, err
	}
	method := abiDefinition.Methods[function.Name]

	transactions, err := s.transactionService.FindDoneTransactionsByFunctionIdAndOrderByTimestamp(functionId, int64(s.cfg.FuzzerConfig.SeedsSize)*2)
	if err != nil {
		return nil, err
	}

	orderer := buildOrderer(strategy, contract)
	orderer.OrderTransactions(transactions)

	seeds := make([][]string, 0)
	for idx := 0; idx < len(transactions) && idx < s.cfg.FuzzerConfig.SeedsSize; idx++ {
		seeds = append(seeds, transactions[idx].Inputs)
	}

	deserializedSeeds, err := s.deserializeSeedsList(method, seeds)
	if err != nil {
		return nil, err
	}

	if len(seeds) < s.cfg.FuzzerConfig.SeedsSize {
		deserializedSeeds, err = s.completeSeedsWithPreConfiguredSeeds(method, deserializedSeeds, s.cfg.FuzzerConfig.SeedsSize-len(seeds))
		if err != nil {
			return nil, err
		}
	}

	return deserializedSeeds, nil
}

func (s *powerSchedule) completeSeedsWithPreConfiguredSeeds(method abi.Method, seeds [][]interface{}, seedsAmountToBeAdded int) ([][]interface{}, error) {
	result := make([][]interface{}, len(seeds)+seedsAmountToBeAdded)
	copy(result, seeds)
	for icr := 0; icr < int(seedsAmountToBeAdded); icr++ {
		functionSeeds := make([]interface{}, len(method.Inputs))
		for inputsIdx, input := range method.Inputs {
			handler, err := s.solidityService.GetTypeHandlerWithContext(input.Type)
			if err != nil {
				return nil, err
			}

			err = handler.LoadSeedsAndChooseOneRandomly(s.cfg.FuzzerConfig.Seeds)
			if err != nil {
				return nil, err
			}

			functionSeeds[inputsIdx] = handler.GetValue()
		}
		result[icr+len(seeds)] = functionSeeds
	}
	return result, nil
}

func (s *powerSchedule) deserializeSeedsList(method abi.Method, seedsList [][]string) ([][]interface{}, error) {
	result := make([][]interface{}, len(seedsList))
	for seedsListIdx, seeds := range seedsList {
		deserializedSeeds := make([]interface{}, len(seeds))
		for inputsIdx, inputDefinition := range method.Inputs {
			if len(seeds) <= inputsIdx {
				return nil, ErrSeedsListInvalid
			}

			handler, err := s.solidityService.GetTypeHandlerWithContext(inputDefinition.Type)
			if err != nil {
				return nil, err
			}

			err = handler.Deserialize(seeds[inputsIdx])
			if err != nil {
				return nil, err
			}
			deserializedSeeds[inputsIdx] = handler.GetValue()
		}
		result[seedsListIdx] = deserializedSeeds
	}
	return result, nil
}
