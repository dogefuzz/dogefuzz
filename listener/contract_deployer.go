package listener

import (
	"context"
	"fmt"
	"strings"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/distance"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"go.uber.org/zap"
)

type contractDeployerListener struct {
	cfg                   *config.Config
	logger                *zap.Logger
	taskStartTopic        interfaces.Topic[bus.TaskStartEvent]
	taskInputRequestTopic interfaces.Topic[bus.TaskInputRequestEvent]
	taskFinishTopic       interfaces.Topic[bus.TaskFinishEvent]
	taskService           interfaces.TaskService
	gethService           interfaces.GethService
	vandalService         interfaces.VandalService
	contractService       interfaces.ContractService
	functionService       interfaces.FunctionService
	transactionService    interfaces.TransactionService
	solidityService       interfaces.SolidityService
	contractMapper        interfaces.ContractMapper
}

func NewContractDeployerListener(e Env) *contractDeployerListener {
	return &contractDeployerListener{
		cfg:                   e.Config(),
		logger:                e.Logger(),
		taskStartTopic:        e.TaskStartTopic(),
		taskInputRequestTopic: e.TaskInputRequestTopic(),
		taskFinishTopic:       e.TaskFinishTopic(),
		taskService:           e.TaskService(),
		gethService:           e.GethService(),
		vandalService:         e.VandalService(),
		contractService:       e.ContractService(),
		functionService:       e.FunctionService(),
		transactionService:    e.TransactionService(),
		solidityService:       e.SolidityService(),
		contractMapper:        e.ContractMapper(),
	}
}

func (l *contractDeployerListener) Name() string { return "contract-deployer" }
func (l *contractDeployerListener) StartListening(ctx context.Context) {
	handler := func(evt bus.TaskStartEvent) { l.processEvent(ctx, evt) }
	l.taskStartTopic.Subscribe(handler)
	<-ctx.Done()
	l.taskStartTopic.Unsubscribe(handler)
}

func (l *contractDeployerListener) processEvent(ctx context.Context, evt bus.TaskStartEvent) {
	l.logger.Debug("processing TaskStartEvent...")

	task, err := l.taskService.Get(evt.TaskId)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving task: %v", err)
		sendTaskEndDueToError(task, l)
		return
	}

	contract, err := l.contractService.FindByTaskId(task.Id)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving contract: %v", err)
		sendTaskEndDueToError(task, l)
		return
	}

	if contract.Status == common.CONTRACT_DEPLOYED {
		l.logger.Sugar().Errorf("the contract %s was already deployed", contract.Id)
		sendTaskEndDueToError(task, l)
		return
	}

	parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when parsing contract ABI definition: %v", err)
		sendTaskEndDueToError(task, l)
		return
	}

	constructor, err := l.functionService.FindConstructorByContractId(contract.Id)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving contract's constructor: %v", err)
		sendTaskEndDueToError(task, l)
		return
	}

	args := make([]interface{}, 0)
	var idx int64
	for idx = 0; idx < constructor.NumberOfArgs; idx++ {
		definition := parsedABI.Constructor.Inputs[idx]

		handler, err := l.solidityService.GetTypeHandlerWithContext(definition.Type)
		if err != nil {
			l.logger.Sugar().Errorf("an error ocurred when parsing args: %v", err)
			sendTaskEndDueToError(task, l)
			return
		}

		if len(task.Arguments) > 0 {
			err = handler.Deserialize(task.Arguments[idx])
			if err != nil {
				l.logger.Sugar().Errorf("an error ocurred when parsing args: %v", err)
				sendTaskEndDueToError(task, l)
				return
			}
		} else {
			handler.Generate()
		}

		args = append(args, handler.GetValue())
	}

	address, tx, err := l.gethService.Deploy(ctx, l.contractMapper.MapDTOToCommon(contract), args...)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when deploying contract: %v", err)
		sendTaskEndDueToError(task, l)
		return
	}
	l.logger.Sugar().Debugf("deploying contract %s at %s", contract.Id, address)
	contract.Address = address

	transactionDTO := &dto.NewTransactionDTO{
		Timestamp:      common.Now(),
		TaskId:         task.Id,
		FunctionId:     constructor.Id,
		Inputs:         task.Arguments,
		Status:         common.TRANSACTION_DONE,
		BlockchainHash: tx,
	}
	_, err = l.transactionService.Create(transactionDTO)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when storing contract's deploy transaction: %v", err)
		sendTaskEndDueToError(task, l)
		return
	}

	cfg, err := l.vandalService.GetCFG(ctx, l.contractMapper.MapDTOToCommon(contract))
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred while getting CFG from vandal service: %v", err)
		sendTaskEndDueToError(task, l)
		return
	}
	contract.CFG = *cfg
	contract.DistanceMap = distance.ComputeDistanceMap(*cfg, l.cfg.FuzzerConfig.CritialInstructions)
	contract.TargetInstructionsFreq = distance.ComputeTargetInstructionsFrequency(*cfg, l.cfg.FuzzerConfig.CritialInstructions)
	contract.Status = common.CONTRACT_DEPLOYED
	err = l.contractService.Update(contract)
	if err != nil {
		l.logger.Sugar().Errorf("an error occurred while updating contract adress: %v", err)
		sendTaskEndDueToError(task, l)
		return
	}
	l.logger.Sugar().Debugf("generated contract's CFG and distance map for contract %s", contract.Id)

	task.Expiration = common.Now().Add(task.Duration)
	err = l.taskService.Update(task)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred while updating task with new expiration: %v", err)
		sendTaskEndDueToError(task, l)
		return
	}
	l.logger.Sugar().Debugf("updating task with new expiration %s", contract.Id)

	l.logger.Info(fmt.Sprintf("requesting new inputs for task %s", task.Id))
	l.taskInputRequestTopic.Publish(bus.TaskInputRequestEvent{TaskId: task.Id})
}

func sendTaskEndDueToError(task *dto.TaskDTO, l *contractDeployerListener) {
	task.Status = common.TASK_DEPLOY_ERROR
	l.taskService.Update(task)
	l.taskFinishTopic.Publish(bus.TaskFinishEvent{TaskId: task.Id})
}
