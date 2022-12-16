package service

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"go.uber.org/zap"
)

type GethService interface {
	Deploy(ctx context.Context, contract *common.Contract, args ...interface{}) (string, error)
}

type gethService struct {
	logger   *zap.Logger
	deployer geth.Deployer
}

func NewGethService(e Env) *gethService {
	return &gethService{
		logger:   e.Logger(),
		deployer: e.Deployer(),
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
