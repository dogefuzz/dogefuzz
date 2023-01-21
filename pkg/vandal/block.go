package vandal

import (
	"regexp"
	"strings"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

const SECTION_DELIMITER = "---"

func NewBlockFromLines(lines []string) common.Block {
	block := common.Block{}
	blockSections := make([][]string, 0)
	section := make([]string, 0)
	for _, line := range lines {
		if line == SECTION_DELIMITER {
			blockSections = append(blockSections, section)
			section = make([]string, 0)
			continue
		}
		section = append(section, line)
	}
	blockSections = append(blockSections, section)

	readHeadSection(blockSections[0], &block)
	readGraphSection(blockSections[1], &block)
	readOpsSection(blockSections[2], &block)
	readStackOpsSection(blockSections[3], &block)
	readFooterSection(blockSections[4], &block)
	return block
}

func readHeadSection(lines []string, block *common.Block) {
	regexBlockPc := regexp.MustCompile("Block (.*)")
	matchBlockPc := regexBlockPc.FindStringSubmatch(lines[0])
	regexBlockRange := regexp.MustCompile(`\[(.*):(.*)\]`)
	matchBlockRange := regexBlockRange.FindStringSubmatch(lines[1])
	block.PC = matchBlockPc[1]
	block.Range = common.BlockRange{
		From: matchBlockRange[1],
		To:   matchBlockRange[2],
	}
}

func readGraphSection(lines []string, block *common.Block) {
	regexPredecessors := regexp.MustCompile(`Predecessors: \[(.*)\]`)
	matchPredecessors := regexPredecessors.FindStringSubmatch(lines[0])
	predecessors := strings.Split(matchPredecessors[1], ",")
	block.Predecessors = make([]string, 0)
	for _, predecessor := range predecessors {
		value := strings.Trim(predecessor, " ")
		if value != "" {
			block.Predecessors = append(block.Predecessors, value)
		}
	}

	regexSuccessors := regexp.MustCompile(`Successors: \[(.*)\]`)
	matchSuccessors := regexSuccessors.FindStringSubmatch(lines[1])
	successors := strings.Split(matchSuccessors[1], ",")
	block.Successors = make([]string, 0)
	for _, successor := range successors {
		value := strings.Trim(successor, " ")
		if value != "" {
			block.Successors = append(block.Successors, value)
		}
	}
}

func readOpsSection(lines []string, block *common.Block) {
	block.InstructionOrder = make([]string, 0)
	block.Instructions = make(map[string]common.Instruction)
	for _, line := range lines {
		spaceIdx := strings.Index(line, " ")
		pc := line[0:spaceIdx]
		block.InstructionOrder = append(block.InstructionOrder, pc)

		opWithArgs := line[spaceIdx+1:]
		opWithArgsList := strings.Split(opWithArgs, " ")
		block.Instructions[pc] = common.Instruction{
			Op:   opWithArgsList[0],
			Args: opWithArgsList[1:],
		}
	}
}

func readStackOpsSection(lines []string, block *common.Block) {
	for _, line := range lines {
		dividerIdx := strings.Index(line, ": ")
		pc := line[0:dividerIdx]
		instruction := block.Instructions[pc]
		instruction.StackOp = line[dividerIdx+2 : len(line)-1]
		block.Instructions[pc] = instruction
	}
}

func readFooterSection(lines []string, block *common.Block) {
	block.EntryStack = readSlicePropertyLine("Entry stack", lines[0])
	block.StackPops = readIntPropertyLine("Stack pops", lines[1])
	block.StackAdditions = readSlicePropertyLine("Stack additions", lines[2])
	block.ExitStack = readSlicePropertyLine("Exit stack", lines[3])
}
