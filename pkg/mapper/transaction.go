package mapper

import (
	"strconv"
	"strings"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type transactionMapper struct{}

func NewTransactionMapper() *transactionMapper {
	return &transactionMapper{}
}

func (m *transactionMapper) MapNewDTOToEntity(n *dto.NewTransactionDTO) *entities.Transaction {
	return &entities.Transaction{
		Timestamp:  n.Timestamp,
		TaskId:     n.TaskId,
		FunctionId: n.FunctionId,
		Inputs:     strings.Join(n.Inputs, ";"),
		Status:     n.Status,
	}
}

func (m *transactionMapper) MapDTOToEntity(c *dto.TransactionDTO) *entities.Transaction {
	coverage := strconv.FormatUint(c.Coverage, 10)
	deltaCoverage := strconv.FormatUint(c.DeltaCoverage, 10)
	deltaMinDistance := strconv.FormatUint(c.DeltaMinDistance, 10)
	criticalInstructionsHits := strconv.FormatUint(c.CriticalInstructionsHits, 10)

	return &entities.Transaction{
		Id:                       c.Id,
		Timestamp:                c.Timestamp,
		BlockchainHash:           c.BlockchainHash,
		TaskId:                   c.TaskId,
		FunctionId:               c.FunctionId,
		Inputs:                   strings.Join(c.Inputs, ";"),
		DetectedWeaknesses:       strings.Join(c.DetectedWeaknesses, ";"),
		ExecutedInstructions:     strings.Join(c.ExecutedInstructions, ";"),
		Coverage:                 coverage,
		DeltaCoverage:            deltaCoverage,
		DeltaMinDistance:         deltaMinDistance,
		CriticalInstructionsHits: criticalInstructionsHits,
		Status:                   c.Status,
	}
}

func (m *transactionMapper) MapEntityToDTO(c *entities.Transaction) *dto.TransactionDTO {
	var inputs []string
	if c.Inputs != "" {
		inputs = strings.Split(c.Inputs, ";")
	}

	var detectedWeaknesses []string
	if c.DetectedWeaknesses != "" {
		detectedWeaknesses = strings.Split(c.DetectedWeaknesses, ";")
	}

	var executedInstructions []string
	if c.ExecutedInstructions != "" {
		executedInstructions = strings.Split(c.ExecutedInstructions, ";")
	}

	var coverage uint64
	if c.Coverage != "" {
		val, err := strconv.ParseUint(c.Coverage, 10, 64)
		if err != nil {
			panic(err)
		}
		coverage = val
	}

	var deltaCoverage uint64
	if c.DeltaCoverage != "" {
		val, err := strconv.ParseUint(c.DeltaCoverage, 10, 64)
		if err != nil {
			panic(err)
		}
		deltaCoverage = val
	}

	var deltaMinDistance uint64
	if c.DeltaMinDistance != "" {
		val, err := strconv.ParseUint(c.DeltaMinDistance, 10, 64)
		if err != nil {
			panic(err)
		}
		deltaMinDistance = val
	}

	var criticalInstructionsHits uint64
	if c.CriticalInstructionsHits != "" {
		val, err := strconv.ParseUint(c.CriticalInstructionsHits, 10, 64)
		if err != nil {
			panic(err)
		}
		criticalInstructionsHits = val
	}

	return &dto.TransactionDTO{
		Id:                       c.Id,
		Timestamp:                c.Timestamp,
		BlockchainHash:           c.BlockchainHash,
		TaskId:                   c.TaskId,
		FunctionId:               c.FunctionId,
		Inputs:                   inputs,
		DetectedWeaknesses:       detectedWeaknesses,
		ExecutedInstructions:     executedInstructions,
		Coverage:                 coverage,
		DeltaCoverage:            deltaCoverage,
		DeltaMinDistance:         deltaMinDistance,
		CriticalInstructionsHits: criticalInstructionsHits,
		Status:                   c.Status,
	}
}
