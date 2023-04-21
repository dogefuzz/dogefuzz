package service

import (
	"errors"
	"sync"
	"time"

	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type transactionService struct {
	transactionRepo   interfaces.TransactionRepo
	transactionMapper interfaces.TransactionMapper
	connection        interfaces.Connection
	mu                sync.Mutex
}

func NewTransactionService(e Env) *transactionService {
	return &transactionService{
		transactionRepo:   e.TransactionRepo(),
		transactionMapper: e.TransactionMapper(),
		connection:        e.DbConnection(),
	}
}

func (s *transactionService) Get(transactionId string) (*dto.TransactionDTO, error) {
	entity, err := s.transactionRepo.Get(s.connection.GetDB(), transactionId)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return s.transactionMapper.MapEntityToDTO(entity), nil
}

func (s *transactionService) Update(transaction *dto.TransactionDTO) error {
	entity := s.transactionMapper.MapDTOToEntity(transaction)
	err := s.transactionRepo.Update(s.connection.GetDB(), entity)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return ErrTransactionNotFound
		}
		return err
	}
	return nil
}

func (s *transactionService) Create(transaction *dto.NewTransactionDTO) (*dto.TransactionDTO, error) {
	entity := s.transactionMapper.MapNewDTOToEntity(transaction)
	err := s.transactionRepo.Create(s.connection.GetDB(), entity)
	if err != nil {
		return nil, err
	}
	return s.transactionMapper.MapEntityToDTO(entity), nil
}

func (s *transactionService) BulkCreate(newTransactions []*dto.NewTransactionDTO) ([]*dto.TransactionDTO, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx := s.connection.GetDB().Begin()
	trasactions := make([]*dto.TransactionDTO, len(newTransactions))
	for idx, dto := range newTransactions {
		entity := s.transactionMapper.MapNewDTOToEntity(dto)
		err := s.transactionRepo.Create(tx, entity)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		trasactions[idx] = s.transactionMapper.MapEntityToDTO(entity)
	}
	tx.Commit()

	return trasactions, nil
}

func (s *transactionService) BulkUpdate(updatedTransactions []*dto.TransactionDTO) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx := s.connection.GetDB().Begin()
	for _, dto := range updatedTransactions {
		entity := s.transactionMapper.MapDTOToEntity(dto)
		err := s.transactionRepo.Update(tx, entity)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, repo.ErrNotExists) {
				return ErrTransactionNotFound
			}
			return err
		}
	}
	tx.Commit()
	return nil
}

func (s *transactionService) FindByHash(hash string) (*dto.TransactionDTO, error) {
	transaction, err := s.transactionRepo.FindByBlockchainHash(s.connection.GetDB(), hash)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return s.transactionMapper.MapEntityToDTO(transaction), nil
}

func (s *transactionService) FindByTaskId(taskId string) ([]*dto.TransactionDTO, error) {
	transactions, err := s.transactionRepo.FindByTaskId(s.connection.GetDB(), taskId)
	if err != nil {
		return nil, err
	}
	transactionDTOs := make([]*dto.TransactionDTO, len(transactions))
	for idx, transaction := range transactions {
		transactionDTOs[idx] = s.transactionMapper.MapEntityToDTO(&transaction)
	}
	return transactionDTOs, nil
}

func (s *transactionService) FindDoneByTaskId(taskId string) ([]*dto.TransactionDTO, error) {
	transactions, err := s.transactionRepo.FindDoneByTaskId(s.connection.GetDB(), taskId)
	if err != nil {
		return nil, err
	}
	transactionDTOs := make([]*dto.TransactionDTO, len(transactions))
	for idx, transaction := range transactions {
		transactionDTOs[idx] = s.transactionMapper.MapEntityToDTO(&transaction)
	}
	return transactionDTOs, nil
}

func (s *transactionService) FindDoneTransactionsByFunctionIdAndOrderByTimestamp(functionId string, limit int64) ([]*dto.TransactionDTO, error) {
	transactions, err := s.transactionRepo.FindDoneTransactionsByFunctionIdAndOrderByTimestamp(s.connection.GetDB(), functionId, limit)
	if err != nil {
		return nil, err
	}
	transactionDTOs := make([]*dto.TransactionDTO, len(transactions))
	for idx, transaction := range transactions {
		transactionDTOs[idx] = s.transactionMapper.MapEntityToDTO(&transaction)
	}
	return transactionDTOs, nil
}

func (s *transactionService) FindRunningAndCreatedBeforeThreshold(dateThreshold time.Time) ([]*dto.TransactionDTO, error) {
	transactions, err := s.transactionRepo.FindRunningAndCreatedBeforeThreshold(s.connection.GetDB(), dateThreshold)
	if err != nil {
		return nil, err
	}
	transactionDTOs := make([]*dto.TransactionDTO, len(transactions))
	for idx, transaction := range transactions {
		transactionDTOs[idx] = s.transactionMapper.MapEntityToDTO(&transaction)
	}
	return transactionDTOs, nil
}
