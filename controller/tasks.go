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
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type tasksController struct {
	logger           *zap.Logger
	taskService      interfaces.TaskService
	contractService  interfaces.ContractService
	functionService  interfaces.FunctionService
	solidityService  interfaces.SolidityService
	taskStartTopic   interfaces.Topic[bus.TaskStartEvent]
	solidityCompiler interfaces.SolidityCompiler
}

func NewTasksController(e Env) *tasksController {
	return &tasksController{
		logger:           e.Logger(),
		taskService:      e.TaskService(),
		contractService:  e.ContractService(),
		functionService:  e.FunctionService(),
		solidityService:  e.SolidityService(),
		taskStartTopic:   e.TaskStartTopic(),
		solidityCompiler: e.SolidityCompiler(),
	}
}

func (ctrl *tasksController) Start(c *gin.Context) {
	var request dto.StartTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("request failed to be parsed: %s", err.Error())})
		return
	}

	var duration time.Duration
	duration, err := time.ParseDuration(request.Duration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("duration failed to be parsed: %s", err.Error())})
		return
	}

	compiledContract, err := ctrl.solidityCompiler.CompileSource(request.ContractName, request.ContractSource)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("contract failed to be compiled: %s", err.Error())})
		return
	}

	parsedABI, err := abi.JSON(strings.NewReader(compiledContract.AbiDefinition))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("contract failed to be parsed: %s", err.Error())})
		return
	}

	payableMethods := make([]abi.Method, 0)
	for _, function := range parsedABI.Methods {
		if !ctrl.isMethodChangingState(function) {
			continue
		}
		payableMethods = append(payableMethods, function)
	}
	if len(payableMethods) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the provide contract doesn't have methods that changes the ontract's state"})
		return
	}

	if len(request.Arguments) > 0 {
		err = ctrl.tryValidateArgs(parsedABI, request.Arguments)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("args failed to be parsed: %s", err.Error())})
			return
		}
	}

	now := common.Now()
	taskDTO := dto.NewTaskDTO{
		Arguments:   request.Arguments,
		StartTime:   now,
		Expiration:  now.Add(duration),
		Detectors:   common.GetUniqueSlice(request.Detectors),
		FuzzingType: request.FuzzingType,
		Status:      common.TASK_RUNNING,
	}
	task, err := ctrl.taskService.Create(&taskDTO)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	contractDTO := dto.NewContractDTO{
		TaskId:             task.Id,
		Source:             request.ContractSource,
		DeploymentBytecode: compiledContract.DeploymentBytecode,
		RuntimeBytecode:    compiledContract.RuntimeBytecode,
		AbiDefinition:      compiledContract.AbiDefinition,
		Name:               compiledContract.Name,
	}
	contract, err := ctrl.contractService.Create(&contractDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	functionDTO := dto.NewFunctionDTO{
		Name:                    parsedABI.Constructor.Name,
		NumberOfArgs:            int64(len(parsedABI.Constructor.Inputs)),
		IsChangingContractState: ctrl.isMethodChangingState(parsedABI.Constructor),
		IsConstructor:           true,
		ContractId:              contract.Id,
	}
	_, err = ctrl.functionService.Create(&functionDTO)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	for _, method := range payableMethods {
		functionDTO := dto.NewFunctionDTO{
			Name:                    method.Name,
			NumberOfArgs:            int64(len(method.Inputs)),
			IsChangingContractState: ctrl.isMethodChangingState(method),
			IsConstructor:           false,
			ContractId:              contract.Id,
		}
		_, err := ctrl.functionService.Create(&functionDTO)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		}
	}

	ctrl.logger.Info(fmt.Sprintf("Requesting fuzzing task %s for %s until %v", task.Id, contract.Name, task.Expiration))
	ctrl.taskStartTopic.Publish(bus.TaskStartEvent{TaskId: task.Id})
	c.JSON(200, dto.StartTaskResponse{TaskId: task.Id})
}

func (ctrl *tasksController) tryValidateArgs(parsedABI abi.ABI, args []string) error {
	if len(args) != len(parsedABI.Constructor.Inputs) {
		return errors.New("invalid number of arguments")
	}

	for idx, arg := range args {
		definition := parsedABI.Constructor.Inputs[idx]

		handler, err := ctrl.solidityService.GetTypeHandlerWithContext(definition.Type)
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

// Reference: https://docs.soliditylang.org/en/v0.8.17/abi-spec.html#json
func (ctrl *tasksController) isMethodChangingState(method abi.Method) bool {
	return method.StateMutability == "nonpayable" || method.StateMutability == "payable"
}
