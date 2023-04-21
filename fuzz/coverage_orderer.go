package fuzz

import (
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type coverageBasedOrderer struct {
}

func newCoverageBasedOrderer() *coverageBasedOrderer {
	return &coverageBasedOrderer{}
}

func (o *coverageBasedOrderer) OrderTransactions(transactions []*dto.TransactionDTO) {
	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].DeltaCoverage < transactions[j].DeltaCoverage
	})
}
