package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TasksController interface {
	Start(c *gin.Context)
}

type tasksController struct {
	logger           *zap.Logger
	taskService      service.TaskService
	oracleService    service.OracleService
	contractService  service.ContractService
	taskStartTopic   topic.Topic[bus.TaskStartEvent]
	solidityCompiler solc.SolidityCompiler
}

func NewTasksController(e Env) *tasksController {
	return &tasksController{
		logger:           e.Logger(),
		taskService:      e.TaskService(),
		oracleService:    e.OracleService(),
		contractService:  e.ContractService(),
		taskStartTopic:   e.TaskStartTopic(),
		solidityCompiler: e.SolidityCompiler(),
	}
}

func (ctrl *tasksController) Start(c *gin.Context) {
	var request dto.StartTaskRequest
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

	taskDTO := dto.NewTaskDTO{
		Contract:   request.Contract,
		Expiration: time.Now().Add(duration),
		Detectors:  request.Detectors,
		Status:     common.TASK_RUNNING,
	}
	task, err := ctrl.taskService.Create(&taskDTO)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	uniqueDetectors := common.GetUniqueSlice(request.Detectors)
	for _, detectorName := range uniqueDetectors {
		oracle, err := ctrl.oracleService.FindByName(detectorName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = ctrl.taskService.AddOracle(task.Id, oracle.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	compiledContract, err := ctrl.solidityCompiler.CompileSource(request.Contract)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	contractDTO := dto.NewContractDTO{
		Source:        request.Contract,
		CompiledCode:  compiledContract.CompiledCode,
		AbiDefinition: compiledContract.AbiDefinition,
		Name:          compiledContract.Name,
	}
	contract, err := ctrl.contractService.Create(&contractDTO)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctrl.taskService.AddContract(task.Id, contract.Id)

	// TODO: store contract's functions in database

	ctrl.logger.Info(fmt.Sprintf("Requesting fuzzing task %s for %s until %v", task.Id, contract.Name, task.Expiration))
	ctrl.taskStartTopic.Publish(bus.TaskStartEvent{TaskId: task.Id})
	c.JSON(200, dto.StartTaskResponse{TaskId: task.Id})
}
