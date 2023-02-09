package coverage

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func ComputeDeltaCoverage(cfg common.CFG, instructionsExecutedInTransaction []string, instructionsExecutedInTask []string) uint64 {
	currentCoverage := ComputeCoverage(cfg, instructionsExecutedInTask)
	mergedInstructions := common.MergeSortedSlices(instructionsExecutedInTask, instructionsExecutedInTransaction)
	newCoverage := ComputeCoverage(cfg, mergedInstructions)
	if newCoverage < currentCoverage {
		return 0
	}
	return newCoverage - currentCoverage
}

func ComputeCoverage(cfg common.CFG, instructions []string) uint64 {
	if len(instructions) == 0 {
		return 0
	}
	edgesPCs := cfg.GetEdgesPCs()
	instructionsSet := common.NewSet[uint64]()
	for _, instruction := range instructions {
		instructionAsUint := common.MustConvertHexadecimalToInt(instruction).Uint64()
		instructionsSet.Add(instructionAsUint)
	}

	coverageCount := 0
	for _, pc := range edgesPCs {
		pcAsUint := common.MustConvertHexadecimalToInt(pc).Uint64()
		if instructionsSet.Has(pcAsUint) {
			coverageCount++
		}
	}
	return uint64(coverageCount)
}
