package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/dogefuzz/dogefuzz/pkg/solidity"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TasksController interface {
	Start(c *gin.Context)
}

type tasksController struct {
	logger           *zap.Logger
	taskService      service.TaskService
	contractService  service.ContractService
	functionService  service.FunctionService
	taskStartTopic   topic.Topic[bus.TaskStartEvent]
	solidityCompiler solc.SolidityCompiler
}

func NewTasksController(e Env) *tasksController {
	return &tasksController{
		logger:           e.Logger(),
		taskService:      e.TaskService(),
		contractService:  e.ContractService(),
		functionService:  e.FunctionService(),
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

	compiledContract, err := ctrl.solidityCompiler.CompileSource(request.Contract)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedABI, err := abi.JSON(strings.NewReader(compiledContract.AbiDefinition))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(request.Arguments) > 0 {
		err = tryValidateArgs(parsedABI, request.Arguments)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	functionDTO := dto.NewFunctionDTO{Name: parsedABI.Constructor.Name, NumberOfArgs: int64(len(parsedABI.Constructor.Inputs))}
	function, err := ctrl.functionService.Create(&functionDTO)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	contractDTO := dto.NewContractDTO{
		Source:        request.Contract,
		CompiledCode:  compiledContract.CompiledCode,
		AbiDefinition: compiledContract.AbiDefinition,
		Name:          compiledContract.Name,
		ConstructorId: function.Id,
	}
	contract, err := ctrl.contractService.Create(&contractDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskDTO := dto.NewTaskDTO{
		ContractId: contract.Id,
		Expiration: time.Now().Add(duration),
		Detectors:  common.GetUniqueSlice(request.Detectors),
		Status:     common.TASK_RUNNING,
	}
	task, err := ctrl.taskService.Create(&taskDTO)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	for _, method := range parsedABI.Methods {
		functionDTO := dto.NewFunctionDTO{Name: method.Name, NumberOfArgs: int64(len(method.Inputs)), Payable: method.Payable}
		_, err := ctrl.functionService.Create(&functionDTO)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		}
	}

	ctrl.logger.Info(fmt.Sprintf("Requesting fuzzing task %s for %s until %v", task.Id, contract.Name, task.Expiration))
	ctrl.taskStartTopic.Publish(bus.TaskStartEvent{TaskId: task.Id})
	c.JSON(200, dto.StartTaskResponse{TaskId: task.Id})
}

func tryValidateArgs(parsedABI abi.ABI, args []string) error {
	if len(args) != len(parsedABI.Constructor.Inputs) {
		return errors.New("invalid number of arguments")
	}

	for idx, arg := range args {
		definition := parsedABI.Constructor.Inputs[idx]

		handler, err := solidity.GetTypeHandler(definition.Type)
		if err != nil {
			return err
		}

		err = handler.Deserialize(arg)
		if err != nil {
			return err
		}
	}
	return nil
}
