package distance

import "github.com/dogefuzz/dogefuzz/pkg/common"

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
	// TODO: Add logic to compute the min distance
	return 0
}
