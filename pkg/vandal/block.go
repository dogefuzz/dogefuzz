package vandal

import (
	"regexp"
	"strings"
)

const SECTION_DELIMITER = "---"

type Block struct {
	PC               string                 `json:"pc"`
	Range            BlockRange             `json:"range"`
	Predecessors     []string               `json:"predecessors"`
	Successors       []string               `json:"successors"`
	EntryStack       []string               `json:"entryStack"`
	StackPops        uint64                 `json:"stackPops"`
	StackAdditions   []string               `json:"stackAdditions"`
	ExitStack        []string               `json:"exitStack"`
	Instructions     map[string]Instruction `json:"instructions"`
	InstructionOrder []string               `json:"instructionOrder"`
}

type BlockRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Instruction struct {
	Op      string `json:"op"`
	StackOp string `json:"stackOp"`
}

func NewBlockFromLines(lines []string) Block {
	block := Block{}
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

func readHeadSection(lines []string, block *Block) {
	regexBlockPc := regexp.MustCompile("Block (.*)")
	matchBlockPc := regexBlockPc.FindStringSubmatch(lines[0])
	regexBlockRange := regexp.MustCompile(`\[(.*):(.*)\]`)
	matchBlockRange := regexBlockRange.FindStringSubmatch(lines[1])
	block.PC = matchBlockPc[1]
	block.Range = BlockRange{
		matchBlockRange[1],
		matchBlockRange[2],
	}
}

func readGraphSection(lines []string, block *Block) {
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

func readOpsSection(lines []string, block *Block) {
	block.InstructionOrder = make([]string, 0)
	block.Instructions = make(map[string]Instruction)
	for _, line := range lines {
		spaceIdx := strings.Index(line, " ")
		pc := line[0:spaceIdx]
		block.InstructionOrder = append(block.InstructionOrder, pc)
		block.Instructions[pc] = Instruction{
			Op: line[spaceIdx+1 : len(line)-1],
		}
	}
}

func readStackOpsSection(lines []string, block *Block) {
	for _, line := range lines {
		dividerIdx := strings.Index(line, ": ")
		pc := line[0:dividerIdx]
		instruction := block.Instructions[pc]
		instruction.StackOp = line[dividerIdx+2 : len(line)-1]
		block.Instructions[pc] = instruction
	}
}

func readFooterSection(lines []string, block *Block) {
	block.EntryStack = readSlicePropertyLine("Entry stack", lines[0])
	block.StackPops = readIntPropertyLine("Stack pops", lines[1])
	block.StackAdditions = readSlicePropertyLine("Stack additions", lines[2])
	block.ExitStack = readSlicePropertyLine("Exit stack", lines[3])
}
