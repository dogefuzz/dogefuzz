package common

type CFG struct {
	Graph        map[string][]string `json:"graph"`
	Blocks       map[string]CFGBlock `json:"blocks"`
	Instructions map[string]string   `json:"instructions"`
}

func (g CFG) GetEdgesPCs() []string {
	edgesPCs := make([]string, len(g.Blocks))
	idx := 0
	for pc := range g.Blocks {
		edgesPCs[idx] = pc
		idx++
	}
	return edgesPCs
}

func (g CFG) GetReverseGraph() map[string][]string {
	reversed := make(map[string][]string)
	for node, edges := range g.Graph {
		for _, edge := range edges {
			reversed[edge] = append(reversed[edge], node)
		}
	}
	return reversed
}

type CFGBlock struct {
	InitialPC       string            `json:"initialPC"`
	Instructions    map[string]string `json:"instructions"`
	InstructionsPCs []string          `json:"instructionsPCs"`
}
