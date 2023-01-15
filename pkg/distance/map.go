package distance

import (
	"math"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func ComputeDistanceMap(cfg common.CFG, targetInstructions []string) common.DistanceMap {
	targetNodes := findNodesContainingTargetInstructions(cfg, targetInstructions)
	reversedCFG := cfg.GetReverseGraph()
	distanceMap := make(common.DistanceMap)
	for cfgNodePC := range cfg.Graph {
		distanceMap[cfgNodePC] = make(map[string]int64)
		for _, targetNode := range targetNodes {
			distanceMap[cfgNodePC][targetNode] = math.MaxInt64
		}
	}

	for _, targetNode := range targetNodes {
		distanceMap[targetNode][targetNode] = 0
		queue := []string{targetNode}
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			for _, edge := range reversedCFG[current] {
				if distanceMap[edge][targetNode] == math.MaxInt64 {
					distanceMap[edge][targetNode] = distanceMap[current][targetNode] + 1
					queue = append(queue, edge)
				}
			}
		}
	}

	return distanceMap
}

func findNodesContainingTargetInstructions(cfg common.CFG, targetInstructions []string) []string {
	targetNodes := make([]string, 0)
	for nodePC, block := range cfg.Blocks {
		for _, instr := range block.Instructions {
			if common.Contains(targetInstructions, instr) {
				targetNodes = append(targetNodes, nodePC)
				break
			}
		}
	}
	return targetNodes
}
