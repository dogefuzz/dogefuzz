package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gongbell/contractfuzzer/db/domain"
	"github.com/gongbell/contractfuzzer/db/repository"
	"github.com/gongbell/contractfuzzer/server/model"
	"go.uber.org/zap"
)

type TransactionAPI interface {
	Create(c *gin.Context)
}

type DefaultTransactionAPI struct {
	logger                *zap.Logger
	transactionRepository repository.TransactionRepository
	contractRepository    repository.ContractRepository
}

func (api DefaultTransactionAPI) Init(
	logger *zap.Logger,
	transactionRepository repository.TransactionRepository,
	contractRepository repository.ContractRepository,
) DefaultTransactionAPI {
	api.logger = logger
	api.transactionRepository = transactionRepository
	return api
}

func (api DefaultTransactionAPI) Create(c *gin.Context) {
	var request model.TransactionCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract, err := api.contractRepository.FindByAddress(request.ContractAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction := domain.Transaction{}
	transaction.TaskId = request.TaskId
	transaction.BlockchainHash = request.BlockchainHash
	transaction.ContractId = contract.Id
	api.transactionRepository.Create(&transaction)

	response := model.TransactionCreateResponse{TransactionId: transaction.Id}
	c.JSON(200, response)
}
