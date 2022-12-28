package listener

import (
	"context"
	"strings"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/distance"
	"github.com/dogefuzz/dogefuzz/pkg/mapper"
	"github.com/dogefuzz/dogefuzz/pkg/solidity"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/dogefuzz/dogefuzz/topic"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"go.uber.org/zap"
)

type contractDeployerListener struct {
	cfg                   *config.Config
	logger                *zap.Logger
	taskStartTopic        topic.Topic[bus.TaskStartEvent]
	taskInputRequestTopic topic.Topic[bus.TaskInputRequestEvent]
	taskService           service.TaskService
	gethService           service.GethService
	vandalService         service.VandalService
	contractService       service.ContractService
	functionService       service.FunctionService
	contractMapper        mapper.ContractMapper
}

func NewContractDeployerListener(e env) *contractDeployerListener {
	return &contractDeployerListener{
		cfg:                   e.Config(),
		logger:                e.Logger(),
		taskStartTopic:        e.TaskStartTopic(),
		taskInputRequestTopic: e.TaskInputRequestTopic(),
		taskService:           e.TaskService(),
		gethService:           e.GethService(),
		vandalService:         e.VandalService(),
		contractService:       e.ContractService(),
		functionService:       e.FunctionService(),
		contractMapper:        e.ContractMapper(),
	}
}

func (l *contractDeployerListener) StartListening() {
	l.taskStartTopic.Subscribe(l.processEvent)
}

func (l *contractDeployerListener) processEvent(evt bus.TaskStartEvent) {
	task, err := l.taskService.Get(evt.TaskId)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving task: %v", err)
		return
	}

	contract, err := l.contractService.Get(task.ContractId)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving contract: %v", err)
		return
	}

	parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when parsing contract ABI definition: %v", err)
		return
	}

	constructor, err := l.functionService.Get(contract.ConstructorId)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving contract's constructor: %v", err)
		return
	}

	args := make([]interface{}, 0)

	var idx int64
	for idx = 0; idx < constructor.NumberOfArgs; idx++ {
		definition := parsedABI.Constructor.Inputs[idx]

		handler, err := solidity.GetTypeHandler(definition.Type)
		if err != nil {
			l.logger.Sugar().Errorf("an error ocurred when parsing args: %v", err)
			return
		}

		if len(task.Arguments) > 0 {
			err = handler.Deserialize(task.Arguments[idx])
			if err != nil {
				l.logger.Sugar().Errorf("an error ocurred when parsing args: %v", err)
				return
			}
		} else {
			handler.Generate()
		}

		args = append(args, handler.GetValue())
	}

	address, err := l.gethService.Deploy(context.Background(), l.contractMapper.ToCommon(contract), args...)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when deploying contract: %v", err)
		return
	}
	contract.Address = address

	cfg, err := l.vandalService.GetCFG(context.Background(), &common.Contract{Name: contract.Name, AbiDefinition: contract.AbiDefinition, CompiledCode: contract.CompiledCode})
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred while getting CFG from vandal service: %v", err)
		return
	}
	contract.CFG = *cfg
	contract.DistanceMap = distance.ComputeDistanceMap(*cfg, l.cfg.FuzzerConfig.CritialInstructions)

	err = l.contractService.Update(contract)
	if err != nil {
		l.logger.Sugar().Errorf("an error occurred while updating contract adress: %v", err)
		return
	}

	l.taskInputRequestTopic.Publish(bus.TaskInputRequestEvent{TaskId: task.Id})
}
