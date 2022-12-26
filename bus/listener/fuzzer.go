package listener

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/dogefuzz/dogefuzz/bus"
	"github.com/dogefuzz/dogefuzz/bus/topic"
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/fuzz"
	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/solidity"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"go.uber.org/zap"
)

type FuzzerListener interface {
	StartListening()
}

type fuzzerListener struct {
	cfg                   *config.Config
	logger                *zap.Logger
	fuzzerLeader          fuzz.FuzzerLeader
	contractMapper        mapper.ContractMapper
	taskInputRequestTopic topic.Topic[bus.TaskInputRequestEvent]
	taskService           service.TaskService
	functionService       service.FunctionService
	contractService       service.ContractService
	gethService           service.GethService
	transactionService    service.TransactionService
}

func NewFuzzerListener(e env) *fuzzerListener {
	return &fuzzerListener{
		cfg:                   e.Config(),
		logger:                e.Logger(),
		fuzzerLeader:          e.FuzzerLeader(),
		contractMapper:        e.ContractMapper(),
		taskInputRequestTopic: e.TaskInputRequestTopic(),
		taskService:           e.TaskService(),
		functionService:       e.FunctionService(),
		contractService:       e.ContractService(),
		gethService:           e.GethService(),
		transactionService:    e.TransactionService(),
	}
}

func (l *fuzzerListener) StartListening() {
	l.taskInputRequestTopic.Subscribe(l.processEvent)
}

func (l *fuzzerListener) processEvent(evt bus.TaskInputRequestEvent) {
	task, err := l.taskService.Get(evt.TaskId)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving task: %v", err)
		return
	}

	if task.Status != common.TASK_RUNNING {
		l.logger.Sugar().Infof("the task %s is not running", task.Id)
		return
	}

	contract, err := l.contractService.Get(task.ContractId)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving contract: %v", err)
		return
	}

	abiDefinition, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving contract's ABI definition: %v", err)
		return
	}

	functions := l.functionService.FindByContractId(task.ContractId)
	chosenFunction := chooseFunction(functions)

	fuzzer, err := l.fuzzerLeader.GetFuzzer(task.FuzzingType)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when getting the fuzzer instance for %s type: %v", task.FuzzingType, err)
		return
	}

	transactionsDTO := make([]*dto.NewTransactionDTO, l.cfg.FuzzerConfig.BatchSize)
	for idx := 0; idx < l.cfg.FuzzerConfig.BatchSize; idx++ {
		inputs := fuzzer.GenerateInput(abiDefinition.Methods[chosenFunction.Name])

		serializedInputs := make([]string, len(inputs))
		abiFunction := abiDefinition.Methods[chosenFunction.Name]
		for idx := 0; idx < len(inputs); idx++ {
			typeHandler, err := solidity.GetTypeHandler(abiFunction.Inputs[idx].Type)
			if err != nil {
				l.logger.Sugar().Errorf("an error ocurred when getting the solidity type handler: %v", err)
				return
			}
			typeHandler.SetValue(inputs[idx])
			serializedInputs[idx] = typeHandler.Serialize()
		}

		transactionsDTO[idx] = &dto.NewTransactionDTO{
			TaskId:     task.Id,
			FunctionId: chosenFunction.Id,
			Inputs:     serializedInputs,
			Status:     common.TRANSACTION_RUNNING,
		}
	}

	transactions, err := l.transactionService.BulkCreate(transactionsDTO)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when creating transactions in database: %v", err)
		return
	}

	inputsByTransactionId := make(map[string][]interface{})
	transactionsByTransactionId := make(map[string]*dto.TransactionDTO)
	for _, tx := range transactions {
		deserializedInputs := make([]interface{}, len(tx.Inputs))
		abiFunction := abiDefinition.Methods[chosenFunction.Name]
		for idx := 0; idx < len(tx.Inputs); idx++ {
			typeHandler, err := solidity.GetTypeHandler(abiFunction.Inputs[idx].Type)
			if err != nil {
				l.logger.Sugar().Errorf("an error ocurred when getting the solidity type handler: %v", err)
				return
			}

			err = typeHandler.Deserialize(tx.Inputs[idx])
			if err != nil {
				l.logger.Sugar().Errorf("an error ocurred when deserialized input: %v", err)
				return
			}
			deserializedInputs[idx] = typeHandler.GetValue()
		}

		inputsByTransactionId[tx.Id] = deserializedInputs
		transactionsByTransactionId[tx.Id] = tx
	}

	transactionHashesByTransactionId, err := l.gethService.BatchCall(context.Background(), l.contractMapper.ToCommon(contract), chosenFunction.Name, inputsByTransactionId)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when sending the transactions: %v", err)
		return
	}

	for transactionId, transactionHash := range transactionHashesByTransactionId {
		transaction := transactionsByTransactionId[transactionId]
		transaction.BlockchainHash = transactionHash
	}

	err = l.transactionService.BulkUpdate(transactions)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when updating transactions in database: %v", err)
		return
	}
}

func chooseFunction(functions []*dto.FunctionDTO) *dto.FunctionDTO {
	rand.Seed(time.Now().Unix())
	idx := rand.Intn(len(functions))
	return functions[idx]
}
