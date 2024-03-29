package distance

import (
	"math"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func ComputeDistanceMap(cfg common.CFG, targetInstructions []string) common.DistanceMap {
	targetBlocks := findBlocksContainingTargetInstructions(cfg, targetInstructions)
	reversedCFG := cfg.GetReverseGraph()
	distanceMap := make(common.DistanceMap)
	for cfgBlockPC := range cfg.Graph {
		distanceMap[cfgBlockPC] = make(map[string]uint32)
		for _, targetBlock := range targetBlocks {
			distanceMap[cfgBlockPC][targetBlock] = math.MaxUint32
		}
	}

	for _, targetBlock := range targetBlocks {
		distanceMap[targetBlock][targetBlock] = 0
		queue := []string{targetBlock}
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			for _, edge := range reversedCFG[current] {
				if distanceMap[edge][targetBlock] == math.MaxUint32 {
					distanceMap[edge][targetBlock] = distanceMap[current][targetBlock] + 1
					queue = append(queue, edge)
				}
			}
		}
	}

	return distanceMap
}

// breath first search over cfg graph searching for target instructions
func ComputeTargetInstructionsFrequency(cfg common.CFG, targetInstructions []string) uint64 {
	var count uint64 = 0
	for _, block := range cfg.Blocks {
		for _, instr := range block.Instructions {
			if common.Contains(targetInstructions, instr) {
				count++
			}
		}
	}
	return count
}

func findBlocksContainingTargetInstructions(cfg common.CFG, targetInstructions []string) []string {
	targetBlocks := make([]string, 0)
	for blockPC, block := range cfg.Blocks {
		for _, instr := range block.Instructions {
			if common.Contains(targetInstructions, instr) {
				targetBlocks = append(targetBlocks, blockPC)
				break
			}
		}
	}
	return targetBlocks
}
