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
		ctrl.logger.Error("request failed to be parsed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("request failed to be parsed: %s", err.Error())})
		return
	}

	var duration time.Duration
	duration, err := time.ParseDuration(request.Duration)
	if err != nil {
		ctrl.logger.Error("duration failed to be parsed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("duration failed to be parsed: %s", err.Error())})
		return
	}

	compiledContract, err := ctrl.solidityCompiler.CompileSource(request.ContractName, request.ContractSource)
	if err != nil {
		ctrl.logger.Error("contract failed to be compiled", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("contract failed to be compiled: %s", err.Error())})
		return
	}

	parsedABI, err := abi.JSON(strings.NewReader(compiledContract.AbiDefinition))
	if err != nil {
		ctrl.logger.Error("contract failed to be parsed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("contract failed to be parsed: %s", err.Error())})
		return
	}

	callableMethods := make([]abi.Method, 0)
	for _, function := range parsedABI.Methods {
		if !isMethodChangingState(function) {
			continue
		}
		callableMethods = append(callableMethods, function)
	}
	if len(callableMethods) == 0 && parsedABI.Fallback.String() == "" && parsedABI.Receive.String() == "" {
		ctrl.logger.Error("the provide contract doesn't have methods that changes the contract's state")
		c.JSON(http.StatusBadRequest, gin.H{"error": "the provide contract doesn't have methods that changes the contract's state"})
		return
	}

	if len(request.Arguments) > 0 {
		err = ctrl.tryValidateArgs(parsedABI, request.Arguments)
		if err != nil {
			ctrl.logger.Error("args failed to be parsed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("args failed to be parsed: %s", err.Error())})
			return
		}
	}

	now := common.Now()
	taskDTO := dto.NewTaskDTO{
		Arguments:   request.Arguments,
		Duration:    duration,
		StartTime:   now,
		Expiration:  now.Add(duration),
		Detectors:   common.GetUniqueSlice(request.Detectors),
		FuzzingType: request.FuzzingType,
		Status:      common.TASK_RUNNING,
	}
	task, err := ctrl.taskService.Create(&taskDTO)
	if err != nil {
		ctrl.logger.Error("task failed to be created", zap.Error(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	contractDTO := dto.NewContractDTO{
		TaskId:             task.Id,
		Status:             common.CONTRACT_CREATED,
		Source:             request.ContractSource,
		DeploymentBytecode: compiledContract.DeploymentBytecode,
		RuntimeBytecode:    compiledContract.RuntimeBytecode,
		AbiDefinition:      compiledContract.AbiDefinition,
		Name:               compiledContract.Name,
	}
	contract, err := ctrl.contractService.Create(&contractDTO)
	if err != nil {
		ctrl.logger.Error("contract failed to be created", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ctrl.storeAvailableMethods(contract, parsedABI)
	if err != nil {
		ctrl.logger.Error("failed to store available methods", zap.Error(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctrl.logger.Info(fmt.Sprintf("Requesting fuzzing task %s for %s until %v", task.Id, contract.Name, task.Expiration))
	ctrl.taskStartTopic.Publish(bus.TaskStartEvent{TaskId: task.Id})
	c.JSON(http.StatusOK, dto.StartTaskResponse{TaskId: task.Id})
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

func (ctrl *tasksController) storeAvailableMethods(contract *dto.ContractDTO, parsedABI abi.ABI) error {
	err := ctrl.storeMethod(contract, parsedABI.Constructor, common.CONSTRUCTOR)
	if err != nil {
		return err
	}

	if parsedABI.Fallback.String() != "" {
		err = ctrl.storeMethod(contract, parsedABI.Fallback, common.FALLBACK)
		if err != nil {
			return err
		}
	}

	if parsedABI.Receive.String() != "" {
		err = ctrl.storeMethod(contract, parsedABI.Receive, common.RECEIVE)
		if err != nil {
			return err
		}
	}

	for _, method := range parsedABI.Methods {
		err = ctrl.storeMethod(contract, method, common.METHOD)
		if err != nil {
			return err
		}
	}
	return nil
}

// stores method in the database
func (ctrl *tasksController) storeMethod(contract *dto.ContractDTO, method abi.Method, methodType common.MethodType) error {
	functionDTO := dto.NewFunctionDTO{
		Name:         method.Name,
		NumberOfArgs: int64(len(method.Inputs)),
		Callable:     isMethodChangingState(method),
		Type:         methodType,
		ContractId:   contract.Id,
	}
	_, err := ctrl.functionService.Create(&functionDTO)
	if err != nil {
		ctrl.logger.Error("function failed to be created", zap.Error(err))
		return err
	}
	return nil
}

// Reference: https://docs.soliditylang.org/en/v0.8.17/abi-spec.html#json
func isMethodChangingState(method abi.Method) bool {
	return method.Payable ||
		method.StateMutability == "nonpayable" ||
		method.StateMutability == "payable"
}
