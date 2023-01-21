package config

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/reporter"
)

type ReporterConfig struct {
	Type            reporter.ReporterType `mapstructure:"type"`
	WebhookEndpoint string                `mapstructure:"webhookEndpoint"`
	WebhookTimeout  time.Duration         `mapstructure:"webhookTimeout"`
}
