package coverage

import "github.com/dogefuzz/dogefuzz/pkg/common"

func ComputeDeltaCoverage(cfg common.CFG, instructionsExecutedInTransaction []string, instructionsExecutedInTask []string) int64 {
	// TODO: Add logic to compute the delta coverage following an exploration strategy
	return 0
}

func ComputeCoverage(cfg common.CFG, instructions []string) int64 {
	// TODO: Add logic to compute the coverage
	return 0
}
