package mapper

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
)

type TransactionMapper interface {
	ToDomain(c *dto.TransactionDTO) *domain.Transaction
	ToDTO(c *domain.Transaction) *dto.TransactionDTO
}

type transactionMapper struct{}

func NewTransactionMapper() *transactionMapper {
	return &transactionMapper{}
}

func (m *transactionMapper) ToDomain(c *dto.TransactionDTO) *domain.Transaction {
	return &domain.Transaction{
		Id:                   c.Id,
		BlockchainHash:       c.BlockchainHash,
		TaskId:               c.TaskId,
		ContractId:           c.ContractId,
		DetectedWeaknesses:   strings.Join(c.DetectedWeaknesses, ";"),
		ExecutedInstructions: strings.Join(c.ExecutedInstructions, ";"),
	}
}

func (m *transactionMapper) ToDTO(c *domain.Transaction) *dto.TransactionDTO {
	return &dto.TransactionDTO{
		Id:                   c.Id,
		BlockchainHash:       c.BlockchainHash,
		TaskId:               c.TaskId,
		ContractId:           c.ContractId,
		DetectedWeaknesses:   strings.Split(c.DetectedWeaknesses, ";"),
		ExecutedInstructions: strings.Split(c.ExecutedInstructions, ";"),
	}
}
