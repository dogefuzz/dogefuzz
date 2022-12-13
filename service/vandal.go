package service

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/vandal"
)

type VandalService interface {
	GetCFG(ctx context.Context, contract *common.Contract) ([]vandal.Block, error)
}

type vandalService struct {
}

func NewVandalService(e Env) *vandalService {
	return &vandalService{}
}

func (s *vandalService) GetCFG(ctx context.Context, contract *common.Contract) ([]vandal.Block, error) {
	// TODO: get logic to get CFG from vandal API
	return nil, nil
}
