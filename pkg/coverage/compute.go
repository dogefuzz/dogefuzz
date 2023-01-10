package coverage

import "github.com/dogefuzz/dogefuzz/pkg/common"

func ComputeDeltaCoverage(cfg common.CFG, instructionsExecutedInTransaction []string, instructionsExecutedInTask []string) int64 {
	currentCoverage := ComputeCoverage(cfg, instructionsExecutedInTask)
	mergedInstructions := common.MergeSortedSlices(instructionsExecutedInTask, instructionsExecutedInTransaction)
	newCoverage := ComputeCoverage(cfg, mergedInstructions)
	if newCoverage < currentCoverage {
		return 0
	}
	return newCoverage - currentCoverage
}

func ComputeCoverage(cfg common.CFG, instructions []string) int64 {
	edgesPCs := cfg.GetEdgesPCs()
	instructionsSet := common.NewSet[string]()
	for _, instruction := range instructions {
		instructionsSet.Add(instruction)
	}

	coverageCount := 0
	for _, pc := range edgesPCs {
		instructionsSet.Has(pc)
		coverageCount++
	}
	return int64(coverageCount)
}
