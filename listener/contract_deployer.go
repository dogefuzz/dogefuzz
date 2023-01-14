package listener

import (
	"context"
	"strings"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/distance"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/pkg/solidity"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"go.uber.org/zap"
)

type contractDeployerListener struct {
	cfg                   *config.Config
	logger                *zap.Logger
	taskStartTopic        interfaces.Topic[bus.TaskStartEvent]
	taskInputRequestTopic interfaces.Topic[bus.TaskInputRequestEvent]
	taskService           interfaces.TaskService
	gethService           interfaces.GethService
	vandalService         interfaces.VandalService
	contractService       interfaces.ContractService
	functionService       interfaces.FunctionService
	contractMapper        interfaces.ContractMapper
}

func NewContractDeployerListener(e Env) *contractDeployerListener {
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

func (l *contractDeployerListener) Name() string { return "contract-deployer" }
func (l *contractDeployerListener) StartListening(ctx context.Context) {
	handler := func(evt bus.TaskStartEvent) { l.processEvent(ctx, evt) }
	l.taskStartTopic.Subscribe(handler)
	<-ctx.Done()
	l.taskStartTopic.Unsubscribe(handler)
}

func (l *contractDeployerListener) processEvent(ctx context.Context, evt bus.TaskStartEvent) {
	task, err := l.taskService.Get(evt.TaskId)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving task: %v", err)
		return
	}

	contract, err := l.contractService.FindByTaskId(task.Id)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving contract: %v", err)
		return
	}

	parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when parsing contract ABI definition: %v", err)
		return
	}

	constructor, err := l.functionService.FindConstructorByContractId(contract.Id)
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

	_, err = l.gethService.Deploy(ctx, l.contractMapper.MapDTOToCommon(contract), args...)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when deploying contract: %v", err)
		return
	}

	cfg, err := l.vandalService.GetCFG(ctx, l.contractMapper.MapDTOToCommon(contract))
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
