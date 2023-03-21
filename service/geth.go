package service

import (
	"context"
	"errors"
	"math/big"
	"math/rand"
	"strings"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/data/repo"
	"github.com/dogefuzz/dogefuzz/environment"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"go.uber.org/zap"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

var ErrInvalidAddress = errors.New("the provided json does not correspond to a address type")

type SendHandler = func(ctx context.Context, wallet *geth.Wallet, contract *common.Contract, functionName string, value *big.Int, args []interface{}) (string, error)

type gethService struct {
	logger         *zap.Logger
	cfg            *config.Config
	deployer       interfaces.Deployer
	agent          interfaces.Agent
	contractRepo   interfaces.ContractRepo
	connection     interfaces.Connection
	contractMapper interfaces.ContractMapper
}

func NewGethService(e Env) *gethService {
	return &gethService{
		logger:         e.Logger(),
		deployer:       e.Deployer(),
		agent:          e.Agent(),
		cfg:            e.Config(),
		contractRepo:   e.ContractRepo(),
		connection:     e.DbConnection(),
		contractMapper: e.ContractMapper(),
	}
}

func (s *gethService) Deploy(ctx context.Context, contract *common.Contract, args ...interface{}) (string, string, error) {
	address, tx, err := s.deployer.Deploy(ctx, contract, args...)
	if err != nil {
		return "", "", err
	}
	s.logger.Sugar().Infof("deploying contract %s at address %s", contract.Name, address)
	return address, tx, nil
}

func (s *gethService) BatchCall(
	ctx context.Context,
	contract *common.Contract,
	functionName string,
	inputsByTransactionId map[string][]interface{},
) (map[string]string, map[string]error) {
	hashesByTransactionId := make(map[string]string)
	errorsByTransactionId := make(map[string]error)

	privateKey := common.RandomChoice([]string{s.cfg.GethConfig.AgentPrivateKeyHex, s.cfg.GethConfig.DeployerPrivateKeyHex})
	wallet, err := geth.NewWalletFromPrivateKeyHex(privateKey)
	if err != nil {
		for transactionId := range inputsByTransactionId {
			errorsByTransactionId[transactionId] = err
		}
		return hashesByTransactionId, errorsByTransactionId
	}

	for transactionId, inputs := range inputsByTransactionId {
		availableSendHandlers := []SendHandler{s.sendToContractDirectly, s.sendToContractViaAgentContract}
		send := common.RandomChoice(availableSendHandlers)

		parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
		if err != nil {
			errorsByTransactionId[transactionId] = err
			continue
		}
		var value *big.Int
		if parsedABI.Methods[functionName].IsPayable() {
			rnd := rand.New(rand.NewSource(common.Now().UnixNano()))
			value = new(big.Int).Rand(rnd, new(big.Int).SetUint64(common.ONE_ETHER))
		} else {
			value = big.NewInt(0)
		}
		hash, err := send(ctx, wallet, contract, functionName, value, inputs)
		if err != nil {
			errorsByTransactionId[transactionId] = err
			continue
		}
		hashesByTransactionId[transactionId] = hash
	}

	return hashesByTransactionId, errorsByTransactionId
}

func (s *gethService) sendToContractDirectly(ctx context.Context, wallet *geth.Wallet, contract *common.Contract, functionName string, value *big.Int, args []interface{}) (string, error) {
	return s.agent.Send(ctx, wallet, contract, functionName, value, args...)
}

func (s *gethService) sendToContractViaAgentContract(ctx context.Context, wallet *geth.Wallet, contract *common.Contract, functionName string, value *big.Int, args []interface{}) (string, error) {
	agentContract, err := s.contractRepo.Find(s.connection.GetDB(), environment.REENTRANCY_AGENT_CONTRACT_ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			return "", ErrContractNotFound
		}
		return "", err
	}

	parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		return "", err
	}

	input, err := parsedABI.Pack(functionName, args...)
	if err != nil {
		return "", err
	}

	contractAddress := gethcommon.HexToAddress(contract.Address)
	if (contractAddress == gethcommon.Address{}) {
		return "", ErrInvalidAddress
	}

	contractDTO := s.contractMapper.MapEntityToDTO(agentContract)
	return s.agent.Send(ctx, wallet, s.contractMapper.MapDTOToCommon(contractDTO), "CallContract", value, contractAddress, input)
}
