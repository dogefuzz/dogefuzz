package controller

import (
	"errors"
	"net/http"

	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/event"
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/repo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExecutionsController interface {
	Create(c *gin.Context)
}

type executionsController struct {
	logger          *zap.Logger
	eventBus        bus.EventBus
	transactionRepo repo.TransactionRepo
}

func NewExecutionsController(e Env) *executionsController {
	return &executionsController{
		logger:          e.Logger(),
		eventBus:        e.EventBus(),
		transactionRepo: e.TransactionRepo(),
	}
}

func (ctrl *executionsController) Create(c *gin.Context) {
	var request dto.NewExecutionDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := ctrl.transactionRepo.FindByBlockchainHash(request.TxHash)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.eventBus.Publish("instrument:execution", event.InstrumentExecutionEvent{
		Input:        request.Input,
		Instructions: request.Instructions,
		Transaction:  *transaction,
	})
	c.AbortWithStatus(200)
}
