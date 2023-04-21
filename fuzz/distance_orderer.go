package fuzz

import (
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type distanceBasedOrderer struct {
}

func newDistanceBasedOrderer() *distanceBasedOrderer {
	return &distanceBasedOrderer{}
}

func (o *distanceBasedOrderer) OrderTransactions(transactions []*dto.TransactionDTO) {
	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].DeltaMinDistance < transactions[j].DeltaCoverage
	})
}
