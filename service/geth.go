package service

import (
	"context"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type gethService struct {
	logger   *zap.Logger
	cfg      *config.Config
	deployer interfaces.Deployer
	agent    interfaces.Agent
}

func NewGethService(e Env) *gethService {
	return &gethService{
		logger:   e.Logger(),
		deployer: e.Deployer(),
		agent:    e.Agent(),
		cfg:      e.Config(),
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
		hash, err := s.agent.Send(ctx, wallet, contract, functionName, inputs...)
		if err != nil {
			errorsByTransactionId[transactionId] = err
			continue
		}
		hashesByTransactionId[transactionId] = hash
	}

	return hashesByTransactionId, errorsByTransactionId
}
