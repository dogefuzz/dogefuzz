package mapper

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type TransactionMapper interface {
	ToDomain(c *dto.TransactionDTO) *entities.Transaction
	ToDTO(c *entities.Transaction) *dto.TransactionDTO
}

type transactionMapper struct{}

func NewTransactionMapper() *transactionMapper {
	return &transactionMapper{}
}

func (m *transactionMapper) ToDomain(c *dto.TransactionDTO) *entities.Transaction {
	return &entities.Transaction{
		Id:                   c.Id,
		BlockchainHash:       c.BlockchainHash,
		TaskId:               c.TaskId,
		FunctionId:           c.FunctionId,
		Inputs:               c.Inputs,
		DetectedWeaknesses:   strings.Join(c.DetectedWeaknesses, ";"),
		ExecutedInstructions: strings.Join(c.ExecutedInstructions, ";"),
	}
}

func (m *transactionMapper) ToDTO(c *entities.Transaction) *dto.TransactionDTO {
	return &dto.TransactionDTO{
		Id:                   c.Id,
		BlockchainHash:       c.BlockchainHash,
		TaskId:               c.TaskId,
		FunctionId:           c.FunctionId,
		Inputs:               c.Inputs,
		DetectedWeaknesses:   strings.Split(c.DetectedWeaknesses, ";"),
		ExecutedInstructions: strings.Split(c.ExecutedInstructions, ";"),
	}
}
