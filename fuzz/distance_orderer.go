package fuzz

import (
	"math"
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

const HITS_WEIGHT = 0.9
const DISTANCE_WEIGHT = 0.1

type distanceBasedOrderer struct {
	contract *dto.ContractDTO
}

func newDistanceBasedOrderer(contract *dto.ContractDTO) *distanceBasedOrderer {
	return &distanceBasedOrderer{contract}
}

func (o *distanceBasedOrderer) OrderTransactions(transactions []*dto.TransactionDTO) {
	sort.SliceStable(transactions, func(i, j int) bool {
		return o.computeScore(transactions[i]) > o.computeScore(transactions[j])
	})
}

func (o *distanceBasedOrderer) computeScore(transaction *dto.TransactionDTO) float64 {
	var hitsPercentage float64
	if o.contract.TargetInstructionsFreq == 0 {
		hitsPercentage = 0
	} else {
		hitsPercentage = float64(transaction.CriticalInstructionsHits) / float64(o.contract.TargetInstructionsFreq)
	}

	var maxDistance map[string]uint32
	for _, distance := range o.contract.DistanceMap {
		if maxDistance != nil {
			maxDistance = make(map[string]uint32, 0)
			for pc := range distance {
				maxDistance[pc] = 0
			}
		}

		for instr := range maxDistance {
			if val, ok := distance[instr]; ok {
				if val != math.MaxUint32 && val > maxDistance[instr] {
					maxDistance[instr] = val
				}
			}
		}
	}

	var distanceSum int64
	for _, distance := range maxDistance {
		distanceSum += int64(distance)
	}
	distancePercentage := (float64(distanceSum) - float64(transaction.DeltaMinDistance)) / float64(distanceSum)

	return HITS_WEIGHT*hitsPercentage + DISTANCE_WEIGHT*distancePercentage
}
