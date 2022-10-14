package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gongbell/contractfuzzer/bus"
	"github.com/gongbell/contractfuzzer/db/domain"
	"github.com/gongbell/contractfuzzer/db/repository"
	"github.com/gongbell/contractfuzzer/pkg/common"
	"github.com/gongbell/contractfuzzer/server/model"
	"go.uber.org/zap"
)

type FuzzAPI interface {
	Start(c *gin.Context)
	Stop(c *gin.Context)
}

type DefaultFuzzAPI struct {
	logger                 *zap.Logger
	eventBus               bus.EventBus
	taskRepository         repository.TaskRepository
	oracleRepository       repository.OracleRepository
	taskOracleRepository   repository.TaskOracleRepository
	contractRepository     repository.ContractRepository
	taskContractRepository repository.TaskContractRepository
}

func (api DefaultFuzzAPI) Init(
	logger *zap.Logger,
	eventBus bus.EventBus,
	taskRepository repository.TaskRepository,
	oracleRepository repository.OracleRepository,
	taskOracleRepository repository.TaskOracleRepository,
	contractRepository repository.ContractRepository,
	taskContractRepository repository.TaskContractRepository,
) DefaultFuzzAPI {
	api.logger = logger
	api.eventBus = eventBus
	api.taskRepository = taskRepository
	api.oracleRepository = oracleRepository
	api.taskOracleRepository = taskOracleRepository
	api.contractRepository = contractRepository
	api.taskContractRepository = taskContractRepository
	return api
}

func (api DefaultFuzzAPI) Start(c *gin.Context) {
	var request model.FuzzStartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var duration time.Duration
	duration, err := time.ParseDuration(request.Duration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := domain.Task{}
	task.Duration = int64(duration)
	api.taskRepository.Create(&task)

	uniqueDetectors := common.GetUniqueSlice(request.Detectors)
	for _, detectorName := range uniqueDetectors {
		oracle, err := api.oracleRepository.FindByName(detectorName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		taskOracle := domain.TaskOracle{}
		taskOracle.OracleId = oracle.Id
		taskOracle.TaskId = task.Id
		api.taskOracleRepository.Create(&taskOracle)
	}

	uniqueContracts := common.GetUniqueSlice(request.Contracts)
	for _, contractName := range uniqueContracts {
		contract, err := api.contractRepository.FindByName(contractName)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		taskContract := domain.TaskContract{}
		taskContract.ContractId = contract.Id
		taskContract.TaskId = task.Id
		api.taskContractRepository.Create(&taskContract)
	}

	api.logger.Info(fmt.Sprintf("Requesting fuzzing task %s for %d seconds for %d contracts", task.Id, task.Duration, len(uniqueContracts)))
	api.eventBus.Publish("task:request", task.Id)
	c.JSON(200, model.FuzzStartResponse{TaskId: task.Id})
}

func (api DefaultFuzzAPI) Stop(c *gin.Context) {
	var taskId = c.Param("task_id")

	task, err := api.taskRepository.Find(taskId)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	api.logger.Info(fmt.Sprintf("Stopping fuzzing task %s", task.Id))
	api.eventBus.Publish(fmt.Sprintf("task:request:%s", task.Id))

	err = api.taskRepository.Delete(taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(200)
}
