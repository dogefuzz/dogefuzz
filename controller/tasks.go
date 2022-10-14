package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/domain"
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/repo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TasksController interface {
	Start(c *gin.Context)
	Stop(c *gin.Context)
}

type tasksController struct {
	logger           *zap.Logger
	eventBus         bus.EventBus
	taskRepo         repo.TaskRepo
	oracleRepo       repo.OracleRepo
	taskOracleRepo   repo.TaskOracleRepo
	contractRepo     repo.ContractRepo
	taskContractRepo repo.TaskContractRepo
}

func NewTasksController(e Env) *tasksController {
	return &tasksController{
		logger:           e.Logger(),
		eventBus:         e.EventBus(),
		taskRepo:         e.TaskRepo(),
		oracleRepo:       e.OracleRepo(),
		taskOracleRepo:   e.TaskOracleRepo(),
		contractRepo:     e.ContractRepo(),
		taskContractRepo: e.TaskContractRepo(),
	}
}

func (ctrl *tasksController) Start(c *gin.Context) {
	var request dto.NewTaskDTO
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
	ctrl.taskRepo.Create(&task)

	uniqueDetectors := common.GetUniqueSlice(request.Detectors)
	for _, detectorName := range uniqueDetectors {
		oracle, err := ctrl.oracleRepo.FindByName(detectorName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		taskOracle := domain.TaskOracle{}
		taskOracle.OracleId = oracle.Id
		taskOracle.TaskId = task.Id
		ctrl.taskOracleRepo.Create(&taskOracle)
	}

	uniqueContracts := common.GetUniqueSlice(request.Contracts)
	for _, contractName := range uniqueContracts {
		contract, err := ctrl.contractRepo.FindByName(contractName)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		taskContract := domain.TaskContract{}
		taskContract.ContractId = contract.Id
		taskContract.TaskId = task.Id
		ctrl.taskContractRepo.Create(&taskContract)
	}

	ctrl.logger.Info(fmt.Sprintf("Requesting fuzzing task %s for %d seconds for %d contracts", task.Id, task.Duration, len(uniqueContracts)))
	ctrl.eventBus.Publish("task:request", task.Id)
	c.JSON(200, dto.TaskDTO{TaskId: task.Id})
}

func (ctrl *tasksController) Stop(c *gin.Context) {
	var taskId = c.Param("task_id")

	task, err := ctrl.taskRepo.Find(taskId)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.logger.Info(fmt.Sprintf("Stopping fuzzing task %s", task.Id))
	ctrl.eventBus.Publish(fmt.Sprintf("task:request:%s", task.Id))

	err = ctrl.taskRepo.Delete(taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(200)
}
