package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/oracle"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TransactionsController interface {
	StoreDetectedWeaknesses(c *gin.Context)
	StoreTransactionExecution(c *gin.Context)
}

type transactionsController struct {
	logger                   *zap.Logger
	transactionService       service.TransactionService
	taskService              service.TaskService
	instrumentExecutionTopic topic.Topic[bus.InstrumentExecutionEvent]
}

func NewTransactionsController(e Env) *transactionsController {
	return &transactionsController{
		logger:                   e.Logger(),
		transactionService:       e.TransactionService(),
		taskService:              e.TaskService(),
		instrumentExecutionTopic: e.InstrumentExecutionTopic(),
	}
}

func (ctrl *transactionsController) StoreDetectedWeaknesses(c *gin.Context) {
	var request dto.NewWeaknessDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := ctrl.transactionService.FindByHash(request.TxHash)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	task, err := ctrl.taskService.Get(transaction.TaskId)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	snapshot := oracle.NewEventsSnapshot(request.OracleEvents)
	var weaknesses []string
	for _, o := range oracle.GetOracles(task.Detectors) {
		if o.Detect(snapshot) {
			weaknesses = append(weaknesses, string(o.Name()))
		}
	}
	transaction.DetectedWeaknesses = weaknesses

	err = ctrl.transactionService.Update(transaction)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(200)
}

func (ctrl *transactionsController) StoreTransactionExecution(c *gin.Context) {
	var request dto.NewExecutionDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := ctrl.transactionService.FindByHash(request.TxHash)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	executedInstructions := make([]string, len(request.Instructions))
	for _, instructionPC := range request.Instructions {
		executedInstructions = append(executedInstructions, strconv.FormatUint(instructionPC, 16))
	}
	transaction.ExecutedInstructions = executedInstructions

	err = ctrl.transactionService.Update(transaction)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.instrumentExecutionTopic.Publish(bus.InstrumentExecutionEvent{TransactionId: transaction.Id})
	c.AbortWithStatus(200)
}
