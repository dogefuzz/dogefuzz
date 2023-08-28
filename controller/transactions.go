package controller

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/pkg/oracle"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ErrTransactionCouldNotBeFoundInDatabase = errors.New("transaction could not be found in database after retries")

type transactionsController struct {
	logger                   *zap.Logger
	transactionService       interfaces.TransactionService
	taskService              interfaces.TaskService
	contractService          interfaces.ContractService
	instrumentExecutionTopic interfaces.Topic[bus.InstrumentExecutionEvent]
	maxRetries               int
}

func NewTransactionsController(e Env) *transactionsController {
	return &transactionsController{
		logger:                   e.Logger(),
		transactionService:       e.TransactionService(),
		taskService:              e.TaskService(),
		contractService:          e.ContractService(),
		instrumentExecutionTopic: e.InstrumentExecutionTopic(),
		maxRetries:               1,
	}
}

func (ctrl *transactionsController) StoreDetectedWeaknesses(c *gin.Context) {
	var request dto.NewWeaknessDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.TxHash == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	transaction, err := ctrl.waitForTransactionToBeStoredInDatabase(request.TxHash)
	if err != nil {
		if errors.Is(err, ErrTransactionCouldNotBeFoundInDatabase) {
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

	if task.Status == common.TASK_DONE {
		ctrl.logger.Info("as the task is done, this weaknesses detection will be ignored")
		c.AbortWithStatus(200)
		return
	}

	snapshot := oracle.NewEventsSnapshot(request.OracleEvents)
	weaknesses := make([]string, 0)
	oracles := oracle.GetOracles(task.Detectors)
	for _, o := range oracles {
		if o.Detect(snapshot) {
			weaknesses = append(weaknesses, string(o.Name()))
		}
	}
	transaction.DetectedWeaknesses = common.MergeSortedSlices(transaction.DetectedWeaknesses, weaknesses)

	err = ctrl.transactionService.Update(transaction)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.logger.Sugar().Infof("storing weaknesses for transaction %s", transaction.Id)
	c.AbortWithStatus(200)
}

func (ctrl *transactionsController) StoreTransactionExecution(c *gin.Context) {
	var request dto.NewExecutionDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.TxHash == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	transaction, err := ctrl.waitForTransactionToBeStoredInDatabase(request.TxHash)
	if err != nil {
		if errors.Is(err, ErrTransactionCouldNotBeFoundInDatabase) {
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

	if task.Status == common.TASK_DONE {
		ctrl.logger.Info("as the task is done, this execution analytics will be ignored")
		c.AbortWithStatus(200)
		return
	}

	contract, err := ctrl.contractService.FindByTaskId(transaction.TaskId)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	if contract.Address != request.Address {
		ctrl.logger.Info("this execution doesn't correspond to the contract being tested, so this execution analytics will be ignored")
		c.AbortWithStatus(200)
		return
	}

	executedInstructions := make([]string, len(request.Instructions))
	for idx, instructionPC := range request.Instructions {
		executedInstructions[idx] = fmt.Sprintf("0x%s", strconv.FormatUint(instructionPC, 16))
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
	ctrl.logger.Sugar().Infof("storing executed instructions for transaction %s", transaction.Id)

	ctrl.instrumentExecutionTopic.Publish(bus.InstrumentExecutionEvent{TransactionId: transaction.Id})
	ctrl.logger.Sugar().Infof("request execution analysis for transaction %s", transaction.Id)
	c.AbortWithStatus(200)
}

func (ctrl *transactionsController) waitForTransactionToBeStoredInDatabase(transactionHash string) (*dto.TransactionDTO, error) {
	var transaction *dto.TransactionDTO
	attemptCounter := 0
	for {
		tx, err := ctrl.transactionService.FindByHash(transactionHash)
		if err != nil {
			if errors.Is(err, service.ErrTransactionNotFound) {
				attemptCounter++
				if attemptCounter >= ctrl.maxRetries {
					return nil, ErrTransactionCouldNotBeFoundInDatabase
				}
				time.Sleep(time.Duration(int(math.Pow(2, float64(attemptCounter)))) * time.Second)
				continue
			}
			return nil, err
		}
		transaction = tx
		break
	}
	return transaction, nil
}
