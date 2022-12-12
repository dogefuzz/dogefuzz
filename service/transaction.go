package service

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/repo"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type TransactionService interface {
	FindByHash(hash string) (*dto.TransactionDTO, error)
	Update(transaction *dto.TransactionDTO) error
}

type transactionService struct {
	transactionRepo   repo.TransactionRepo
	transactionMapper mapper.TransactionMapper
}

func NewTransactionService(e Env) *transactionService {
	return &transactionService{
		transactionRepo:   e.TransactionRepo(),
		transactionMapper: e.TransactionMapper(),
	}
}

func (s *transactionService) FindByHash(hash string) (*dto.TransactionDTO, error) {
	transaction, err := s.transactionRepo.FindByBlockchainHash(hash)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return s.transactionMapper.ToDTO(transaction), nil
}

func (s *transactionService) Update(transaction *dto.TransactionDTO) error {
	entity := s.transactionMapper.ToDomain(transaction)
	err := s.transactionRepo.Update(entity)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return ErrTransactionNotFound
		}
		return err
	}
	return nil
}
