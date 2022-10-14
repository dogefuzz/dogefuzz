package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gongbell/contractfuzzer/bus"
	"github.com/gongbell/contractfuzzer/bus/event"
	"github.com/gongbell/contractfuzzer/db/repository"
	"github.com/gongbell/contractfuzzer/pkg/oracle"
	"github.com/gongbell/contractfuzzer/server/model"
	"go.uber.org/zap"
)

type InstrumentAPI interface {
	Execution(c *gin.Context)
	Weakness(c *gin.Context)
}

type DefaultInstrumentAPI struct {
	logger                *zap.Logger
	eventBus              bus.EventBus
	transactionRepository repository.TransactionRepository
	taskRepository        repository.TaskRepository
	taskOracleRepository  repository.TaskOracleRepository
	oracleRepository      repository.OracleRepository
}

func (api DefaultInstrumentAPI) Init(
	logger *zap.Logger,
	eventBus bus.EventBus,
	transactionRepository repository.TransactionRepository,
	taskRepository repository.TaskRepository,
	taskOracleRepository repository.TaskOracleRepository,
	oracleRepository repository.OracleRepository,
) DefaultInstrumentAPI {
	api.logger = logger
	api.eventBus = eventBus
	api.transactionRepository = transactionRepository
	api.taskRepository = taskRepository
	api.taskOracleRepository = taskOracleRepository
	api.oracleRepository = oracleRepository
	return api
}

func (api DefaultInstrumentAPI) Execution(c *gin.Context) {
	var request model.InstrumentExecutionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := api.transactionRepository.FindByBlockchainHash(request.TxHash)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	evt := event.InstrumentExecutionEvent{}
	evt.Input = request.Input
	evt.Instructions = request.Instructions
	evt.Transaction = *transaction
	api.eventBus.Publish("instrument:execution", evt)
	c.AbortWithStatus(200)
}

func (api DefaultInstrumentAPI) Weakness(c *gin.Context) {
	var request model.InstrumentWeaknessRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := api.transactionRepository.FindByBlockchainHash(request.TxHash)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	task, err := api.taskRepository.Find(transaction.TaskId)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tasksOracles, err := api.taskOracleRepository.FindByTaskId(task.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var oracleIds []string
	for _, to := range tasksOracles {
		oracleIds = append(oracleIds, to.OracleId)
	}

	oracleEntities, err := api.oracleRepository.FindByIds(oracleIds)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var oracleNames []string
	for _, entity := range oracleEntities {
		oracleNames = append(oracleNames, entity.Name)
	}
	oracles := oracle.GetOracles(oracleNames)
	snapshot := oracle.NewEventsSnapshot(request.OracleEvents)

	var weaknesses []string
	for _, o := range oracles {
		if o.Detect(snapshot) {
			weaknesses = append(weaknesses, o.Name())
		}
	}

	transaction.SetDetectedWeaknesses(weaknesses)
	api.transactionRepository.Update(transaction)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(200)
}
