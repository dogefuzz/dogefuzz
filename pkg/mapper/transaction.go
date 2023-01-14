package mapper

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type TransactionMapper interface {
	MapNewDTOToEntity(n *dto.NewTransactionDTO) *entities.Transaction
	MapDTOToEntity(c *dto.TransactionDTO) *entities.Transaction
	MapEntityToDTO(c *entities.Transaction) *dto.TransactionDTO
}

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
	return &entities.Transaction{
		Id:                   c.Id,
		Timestamp:            c.Timestamp,
		BlockchainHash:       c.BlockchainHash,
		TaskId:               c.TaskId,
		FunctionId:           c.FunctionId,
		Inputs:               strings.Join(c.Inputs, ";"),
		DetectedWeaknesses:   strings.Join(c.DetectedWeaknesses, ";"),
		ExecutedInstructions: strings.Join(c.ExecutedInstructions, ";"),
		DeltaCoverage:        c.DeltaCoverage,
		DeltaMinDistance:     c.DeltaMinDistance,
		Status:               c.Status,
	}
}

func (m *transactionMapper) MapEntityToDTO(c *entities.Transaction) *dto.TransactionDTO {
	return &dto.TransactionDTO{
		Id:                   c.Id,
		Timestamp:            c.Timestamp,
		BlockchainHash:       c.BlockchainHash,
		TaskId:               c.TaskId,
		FunctionId:           c.FunctionId,
		Inputs:               strings.Split(c.Inputs, ";"),
		DetectedWeaknesses:   strings.Split(c.DetectedWeaknesses, ";"),
		ExecutedInstructions: strings.Split(c.ExecutedInstructions, ";"),
		DeltaCoverage:        c.DeltaCoverage,
		DeltaMinDistance:     c.DeltaMinDistance,
		Status:               c.Status,
	}
}
