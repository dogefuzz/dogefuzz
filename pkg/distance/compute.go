package distance

import (
	"math"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func ComputeDeltaMinDistance(distanceMap common.DistanceMap, instructionsExecutedInTransaction []string, instructionsExecutedInTask []string) uint64 {
	currentMinDistance := ComputeMinDistance(distanceMap, instructionsExecutedInTask)
	mergedInstructions := common.MergeSortedSlices(instructionsExecutedInTask, instructionsExecutedInTransaction)
	newMinDistance := ComputeMinDistance(distanceMap, mergedInstructions)
	if currentMinDistance == math.MaxUint64 {
		return newMinDistance
	}

	if newMinDistance > currentMinDistance {
		return 0
	}
	return currentMinDistance - newMinDistance
}

func ComputeMinDistance(distanceMap common.DistanceMap, instructions []string) uint64 {
	if len(instructions) == 0 {
		return math.MaxUint64
	}
	executedBlocks := findExecutedBlocks(distanceMap, instructions)
	minDistances := computeMinDistancesFromExecutedBlocks(distanceMap, executedBlocks)

	var sum uint64 = 0
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

func computeMinDistancesFromExecutedBlocks(distanceMap common.DistanceMap, executedBlocks []string) map[string]uint64 {
	minDistances := make(map[string]uint64)
	for _, block := range executedBlocks {
		for pc, distance := range distanceMap[block] {
			if _, ok := minDistances[pc]; ok {
				minDistances[pc] = uint64(math.Min(float64(minDistances[pc]), float64(distance)))
			} else {
				minDistances[pc] = distance
			}
		}
	}
	return minDistances
}
