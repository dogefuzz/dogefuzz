package controller

import (
	"net/http"

	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/repo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TransactionsController interface {
	Create(c *gin.Context)
}

type transactionsController struct {
	logger          *zap.Logger
	transactionRepo repo.TransactionRepo
	contractRepo    repo.ContractRepo
}

func NewTransactionsController(e Env) *transactionsController {
	return &transactionsController{
		logger:          e.Logger(),
		transactionRepo: e.TransactionRepo(),
	}
}

func (ctrl transactionsController) Create(c *gin.Context) {
	var request dto.NewTransactionDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract, err := ctrl.contractRepo.FindByAddress(request.ContractAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction := domain.Transaction{}
	transaction.TaskId = request.TaskId
	transaction.BlockchainHash = request.BlockchainHash
	transaction.ContractId = contract.Id
	ctrl.transactionRepo.Create(&transaction)

	response := dto.TransactionDTO{TransactionId: transaction.Id}
	c.JSON(200, response)
}
