package interfaces

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type FuzzerLeader interface {
	GetFuzzerStrategy(typ common.FuzzingType) (Fuzzer, error)
}

type PowerSchedule interface {
	RequestSeeds(functionId string, strategy common.PowerScheduleStrategy) ([][]interface{}, error)
}

type Fuzzer interface {
	GenerateInput(functionId string) ([]interface{}, error)
}
