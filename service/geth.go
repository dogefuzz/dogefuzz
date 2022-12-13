package service

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
)

type GethService interface {
	Deploy(ctx context.Context, contract *common.Contract, args ...string) (string, error)
}

type gethService struct {
	deployer geth.Deployer
}

func NewGethService(e Env) *gethService {
	return &gethService{deployer: e.Deployer()}
}

func (s *gethService) Deploy(ctx context.Context, contract *common.Contract, args ...string) (string, error) {
	address, err := s.deployer.Deploy(ctx, contract, common.ConvertStringArrayToInterfaceArray(args)...)
	if err != nil {
		return "", err
	}
	return address, nil
}
