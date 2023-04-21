package fuzz

import (
	"fmt"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

type env interface {
	Config() *config.Config

	BlackboxFuzzer() interfaces.Fuzzer
	GreyboxFuzzer() interfaces.Fuzzer
	DirectedGreyboxFuzzer() interfaces.Fuzzer
	PowerSchedule() interfaces.PowerSchedule

	TransactionService() interfaces.TransactionService
	SolidityService() interfaces.SolidityService
}

type Orderer interface {
	OrderTransactions(transactions []*dto.TransactionDTO)
}

func buildOrderer(strategy common.PowerScheduleStrategy) Orderer {
	switch strategy {
	case common.COVERAGE_BASED_STRATEGY:
		return newCoverageBasedOrderer()
	case common.DISTANCE_BASED_STRATEGY:
		return newDistanceBasedOrderer()
	default:
		panic(fmt.Sprintf("invalid power schedule strategy: %s", strategy))
	}
}
