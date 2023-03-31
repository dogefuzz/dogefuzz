package service

import (
	"context"
	"errors"
	"os"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/pkg/reporter"
)

var ErrReporterNotImplemented = errors.New("this reporter is not implemented")

type reporterService struct {
	client interfaces.HttpClient
	cfg    *config.Config
}

func NewReporterService(e Env) *reporterService {
	return &reporterService{
		client: e.Client(),
		cfg:    e.Config(),
	}
}

func (s *reporterService) SendReport(ctx context.Context, report common.TaskReport) error {
	reporter, err := s.getReporter(s.cfg.ReporterConfig)
	if err != nil {
		return err
	}

	return reporter.SendOutput(ctx, report)
}

func (s *reporterService) getReporter(reporterConfig config.ReporterConfig) (interfaces.Reporter, error) {
	switch reporterConfig.Type {
	case common.CONSOLE_REPORTER:
		return reporter.NewConsoleReporter(os.Stdout), nil
	case common.WEBHOOK_REPOTER:
		return reporter.NewWebhookReporter(s.client, reporterConfig.WebhookEndpoint, reporterConfig.WebhookTimeout), nil
	case common.FILE_REPORTER:
		return reporter.NewFileReporter(reporterConfig.FileOutputPath), nil
	default:
		return nil, ErrReporterNotImplemented
	}
}
