package listener

import (
	"context"

	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/service"
	"go.uber.org/zap"
)

type ContractDeployerListener interface {
	StartListening()
}

type contractDeployerListener struct {
	logger          *zap.Logger
	taskStartTopic  topic.Topic[bus.TaskStartEvent]
	taskService     service.TaskService
	gethService     service.GethService
	vandalService   service.VandalService
	contractService service.ContractService
}

func NewContractDeployerListener(e env) *contractDeployerListener {
	return &contractDeployerListener{
		logger:          e.Logger(),
		taskStartTopic:  e.TaskStartTopic(),
		taskService:     e.TaskService(),
		gethService:     e.GethService(),
		vandalService:   e.VandalService(),
		contractService: e.ContractService(),
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
	err = l.contractService.Update(contract)
	if err != nil {
		l.logger.Sugar().Errorf("and error occurred while updating contract adress: %v", err)
		return
	}

	// TODO: Call Vandal service

	// TODO: Emit task_input_requests
}
