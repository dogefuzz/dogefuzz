package distance

import (
	"math"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func ComputeDeltaMinDistance(distanceMap common.DistanceMap, instructionsExecutedInTransaction []string, instructionsExecutedInTask []string) int64 {
	currentMinDistance := ComputeMinDistance(distanceMap, instructionsExecutedInTask)
	mergedInstructions := common.MergeSortedSlices(instructionsExecutedInTask, instructionsExecutedInTransaction)
	newMinDistance := ComputeMinDistance(distanceMap, mergedInstructions)
	if newMinDistance > currentMinDistance {
		return 0
	}
	return currentMinDistance - newMinDistance
}

func ComputeMinDistance(distanceMap common.DistanceMap, instructions []string) int64 {
	executedBlocks := findExecutedBlocks(distanceMap, instructions)
	minDistances := computeMinDistancesFromExecutedBlocks(distanceMap, executedBlocks)

	var sum int64 = 0
	for _, distance := range minDistances {
		sum += distance
	}
	return sum
}

func findExecutedBlocks(distanceMap common.DistanceMap, instructions []string) []string {
	blockPCs := make([]string, 0)
	for pc := range distanceMap {
		blockPCs = append(blockPCs, pc)
	}

	executedBlockPcs := make([]string, 0)
	for _, instr := range instructions {
		if common.Contains(blockPCs, instr) {
			executedBlockPcs = append(executedBlockPcs, instr)
		}
	}
	return executedBlockPcs
}

func computeMinDistancesFromExecutedBlocks(distanceMap common.DistanceMap, executedBlock []string) map[string]int64 {
	minDistances := make(map[string]int64)
	for _, block := range executedBlock {
		for pc, distance := range distanceMap[block] {
			if _, ok := minDistances[pc]; ok {
				minDistances[pc] = int64(math.Max(float64(minDistances[pc]), float64(distance)))
			} else {
				minDistances[pc] = distance
			}
		}
	}
	return minDistances
}
