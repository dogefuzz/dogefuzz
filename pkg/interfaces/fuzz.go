package interfaces

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type FuzzerLeader interface {
	GetFuzzerStrategy(typ common.FuzzingType) (Fuzzer, error)
}

type PowerSchedule interface {
	RequestSeeds(method abi.Method, strategy common.PowerScheduleStrategy) ([][]interface{}, error)
}

type Fuzzer interface {
	GenerateInput(method abi.Method) ([]interface{}, error)
}
