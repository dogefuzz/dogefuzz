package mapper

import (
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
	coverage := c.Coverage
	deltaCoverage := c.DeltaCoverage
	deltaMinDistance := c.DeltaMinDistance
	criticalInstructionsHits := c.CriticalInstructionsHits

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

	var coverage uint64 = c.Coverage

	var deltaCoverage uint64 = c.DeltaCoverage

	var deltaMinDistance uint64 = c.DeltaMinDistance

	var criticalInstructionsHits uint64 = c.CriticalInstructionsHits

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
