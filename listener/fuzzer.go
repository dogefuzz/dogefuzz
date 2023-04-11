package listener

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"go.uber.org/zap"
)

type fuzzerListener struct {
	cfg                   *config.Config
	logger                *zap.Logger
	fuzzerLeader          interfaces.FuzzerLeader
	contractMapper        interfaces.ContractMapper
	taskInputRequestTopic interfaces.Topic[bus.TaskInputRequestEvent]
	taskService           interfaces.TaskService
	functionService       interfaces.FunctionService
	contractService       interfaces.ContractService
	gethService           interfaces.GethService
	transactionService    interfaces.TransactionService
	solidityService       interfaces.SolidityService
}

func NewFuzzerListener(e Env) *fuzzerListener {
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
		solidityService:       e.SolidityService(),
	}
}

func (l *fuzzerListener) Name() string { return "fuzzer" }
func (l *fuzzerListener) StartListening(ctx context.Context) {
	handler := func(evt bus.TaskInputRequestEvent) { l.processEvent(ctx, evt) }
	l.taskInputRequestTopic.Subscribe(handler)
	<-ctx.Done()
	l.taskInputRequestTopic.Unsubscribe(handler)
}

func (l *fuzzerListener) processEvent(ctx context.Context, evt bus.TaskInputRequestEvent) {
	l.logger.Debug("processing TaskInputRequestEvent...")

	task, err := l.taskService.Get(evt.TaskId)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving task: %v", err)
		return
	}

	if task.Status != common.TASK_RUNNING {
		l.logger.Sugar().Infof("the task %s is not running", task.Id)
		return
	}

	contract, err := l.contractService.FindByTaskId(task.Id)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving contract: %v", err)
		return
	}

	abiDefinition, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when retrieving contract's ABI definition: %v", err)
		return
	}

	functions, err := l.functionService.FindByContractId(contract.Id)
	if err != nil {
		l.logger.Sugar().Errorf("an error occurred when retrieving contract's functions> %v", err)
		return
	}
	chosenFunction := chooseFunction(functions)

	fuzzer, err := l.fuzzerLeader.GetFuzzerStrategy(task.FuzzingType)
	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when getting the fuzzer instance for %s type: %v", task.FuzzingType, err)
		return
	}

	abiFunction := abiDefinition.Methods[chosenFunction.Name]
	transactionsDTO := make([]*dto.NewTransactionDTO, l.cfg.FuzzerConfig.BatchSize)
	for idx := 0; idx < l.cfg.FuzzerConfig.BatchSize; idx++ {
		inputs, err := fuzzer.GenerateInput(abiFunction)
		if err != nil {
			l.logger.Sugar().Errorf("an error ocurred when generating inputs: %v", err)
			return
		}

		serializedInputs := make([]string, len(inputs))
		for idx := 0; idx < len(inputs); idx++ {
			typeHandler, err := l.solidityService.GetTypeHandlerWithContext(abiFunction.Inputs[idx].Type)
			if err != nil {
				l.logger.Sugar().Errorf("an error ocurred when getting the solidity type handler: %v", err)
				return
			}
			typeHandler.SetValue(inputs[idx])
			serializedInputs[idx] = typeHandler.Serialize()
		}

		transactionsDTO[idx] = &dto.NewTransactionDTO{
			Timestamp:  common.Now(),
			TaskId:     task.Id,
			FunctionId: chosenFunction.Id,
			Inputs:     serializedInputs,
			Status:     common.TRANSACTION_CREATED,
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
		for idx := 0; idx < len(tx.Inputs); idx++ {
			typeHandler, err := l.solidityService.GetTypeHandlerWithContext(abiFunction.Inputs[idx].Type)
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

	transactionHashesByTransactionId, errorsByTransactionId := l.gethService.BatchCall(ctx, l.contractMapper.MapDTOToCommon(contract), chosenFunction.Name, inputsByTransactionId)
	for transactionId, err := range errorsByTransactionId {
		l.logger.Sugar().Warnf("failed to send transaction %s to node: %v", transactionId, err)
		transaction := transactionsByTransactionId[transactionId]
		transaction.Status = common.TRANSACTION_SEND_ERROR
	}

	for transactionId, transactionHash := range transactionHashesByTransactionId {
		transaction := transactionsByTransactionId[transactionId]
		transaction.BlockchainHash = transactionHash
		transaction.Status = common.TRANSACTION_RUNNING
	}

	tries := 0
	for tries <= 5 {
		err = l.transactionService.BulkUpdate(transactions)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}
		tries++
	}

	if err != nil {
		l.logger.Sugar().Errorf("an error ocurred when updating transactions in database: %v", err)
		return
	}
}

func chooseFunction(functions []*dto.FunctionDTO) *dto.FunctionDTO {
	rand.Seed(common.Now().Unix())

	filteredFunctions := make([]*dto.FunctionDTO, 0)
	for _, function := range functions {
		if !function.IsConstructor {
			filteredFunctions = append(filteredFunctions, function)
		}
	}

	idx := rand.Intn(len(filteredFunctions))
	return filteredFunctions[idx]
}
