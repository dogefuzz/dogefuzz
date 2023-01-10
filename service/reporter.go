package service

import (
	"context"
	"errors"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/reporter"
)

var ErrReporterNotImplemented = errors.New("this reporter is not implemented")

type ReporterService interface {
	SendReport(ctx context.Context, report common.TaskReport) error
}

type reporterService struct {
	cfg *config.Config
}

func NewReporterService(e Env) *reporterService {
	return &reporterService{cfg: e.Config()}
}

func (s *reporterService) SendReport(ctx context.Context, report common.TaskReport) error {
	reporter, err := s.getReporter(s.cfg.ReporterConfig)
	if err != nil {
		return err
	}

	return reporter.SendOutput(ctx, report)
}

func (s *reporterService) getReporter(reporterConfig config.ReporterConfig) (reporter.Reporter, error) {
	switch reporterConfig.Type {
	case reporter.CONSOLE_REPORTER:
		return reporter.NewConsoleReporter(), nil
	case reporter.CALLBACK_REPOTER:
		return reporter.NewCallbackReporter(reporterConfig.CallbackEndpoint), nil
	default:
		return nil, ErrReporterNotImplemented
	}
}
