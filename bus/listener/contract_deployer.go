package listener

import (
	"context"

	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/distance"
	"github.com/dogefuzz/dogefuzz/service"
	"go.uber.org/zap"
)

type ContractDeployerListener interface {
	StartListening()
}

type contractDeployerListener struct {
	cfg                   *config.Config
	logger                *zap.Logger
	taskStartTopic        topic.Topic[bus.TaskStartEvent]
	taskInputRequestTopic topic.Topic[bus.TaskInputRequestEvent]
	taskService           service.TaskService
	gethService           service.GethService
	vandalService         service.VandalService
	contractService       service.ContractService
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

	address, err := l.gethService.Deploy(
		context.Background(),
		&common.Contract{Name: contract.Name, AbiDefinition: contract.AbiDefinition, CompiledCode: contract.CompiledCode},
	)
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
	contract.DistanceMap = distance.ComputeDistanceMap(*cfg, l.cfg.CritialInstructions)

	err = l.contractService.Update(contract)
	if err != nil {
		l.logger.Sugar().Errorf("an error occurred while updating contract adress: %v", err)
		return
	}

	l.taskInputRequestTopic.Publish(bus.TaskInputRequestEvent{TaskId: task.Id})
}
