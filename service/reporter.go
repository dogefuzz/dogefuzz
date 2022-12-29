package service

import (
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type ReporterService interface {
	SendReport(report common.TaskReport) error
}

type reporterService struct {
	cfg *config.Config
}

func NewReporterService(e Env) *reporterService {
	return &reporterService{
		cfg: e.Config(),
	}
}

func (s *reporterService) SendReport(report common.TaskReport) error {
	// TODO: sends reports
	return nil
}
