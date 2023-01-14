package service

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type gethService struct {
	logger   *zap.Logger
	deployer interfaces.Deployer
	agent    interfaces.Agent
}

func NewGethService(e Env) *gethService {
	return &gethService{
		logger:   e.Logger(),
		deployer: e.Deployer(),
		agent:    e.Agent(),
	}
}

func (s *gethService) Deploy(ctx context.Context, contract *common.Contract, args ...interface{}) (string, error) {
	address, err := s.deployer.Deploy(ctx, contract, args...)
	if err != nil {
		return "", err
	}
	s.logger.Sugar().Infof("deploying contract %s at address %s", contract.Name, address)
	return address, nil
}

func (s *gethService) BatchCall(
	ctx context.Context,
	contract *common.Contract,
	functionName string,
	inputsByTransactionId map[string][]interface{},
) (map[string]string, map[string]error) {
	hashesByTransactionId := make(map[string]string)
	errorsByTransactionId := make(map[string]error)
	for transactionId, inputs := range inputsByTransactionId {
		hash, err := s.agent.Send(ctx, contract, functionName, inputs)
		if err != nil {
			errorsByTransactionId[transactionId] = err
			continue
		}
		hashesByTransactionId[transactionId] = hash
	}

	return hashesByTransactionId, errorsByTransactionId
}
