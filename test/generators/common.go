package generators

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func CommonContractGen() *common.Contract {
	return &common.Contract{
		Name:               gofakeit.Word(),
		AbiDefinition:      gofakeit.Word(),
		DeploymentBytecode: gofakeit.HexUint256(),
	}
}

func CFGGen() common.CFG {
	rand.Seed(time.Now().UnixNano())
	graphLength := rand.Intn(255)
	graph := make(map[string][]string)
	blocks := make(map[string]common.CFGBlock)

	pcs := make([]string, 0)
	for i := 0; i < graphLength; i++ {
		pc := gofakeit.HexUint64()
		graph[pc] = common.RandomSubSlice(pcs)
		blocks[pc] = CFGBlockGen()
		pcs = append(pcs, pc)
	}

	rangeLength := rand.Intn(255)
	instructions := make(map[string]string)
	for idx := 0; idx < rangeLength; idx++ {
		instruction := gofakeit.HexUint64()
		instructions[instruction] = gofakeit.LetterN(32)
	}

	return common.CFG{Graph: graph, Blocks: blocks, Instructions: instructions}
}

func CFGBlockGen() common.CFGBlock {
	rand.Seed(time.Now().UnixNano())
	rangeLength := rand.Intn(255)

	instructions := make(map[string]string)
	instructionsPCs := make([]string, rangeLength)
	for idx := 0; idx < rangeLength; idx++ {
		instruction := gofakeit.HexUint64()
		instructions[instruction] = gofakeit.LetterN(32)
		instructionsPCs[idx] = instruction
	}

	return common.CFGBlock{
		InitialPC:       gofakeit.HexUint256(),
		Instructions:    instructions,
		InstructionsPCs: instructionsPCs,
	}
}

func DistanceMapGen() common.DistanceMap {
	rand.Seed(time.Now().UnixNano())
	rangeLength := rand.Intn(255)
	targetInstrCount := rand.Intn(255)

	targetInstr := make([]string, targetInstrCount)
	for idx := 0; idx < targetInstrCount; idx++ {
		targetInstr[idx] = gofakeit.HexUint64()
	}

	distanceMap := make(map[string]map[string]uint32)
	for icr := 0; icr < rangeLength; icr++ {
		instruction := gofakeit.HexUint64()

		localDistance := make(map[string]uint32)
		for idx := 0; idx < targetInstrCount; idx++ {
			localDistance[targetInstr[idx]] = gofakeit.Uint32()
		}
		distanceMap[instruction] = localDistance
	}
	return common.DistanceMap(distanceMap)
}
