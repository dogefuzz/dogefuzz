package config

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type ReporterConfig struct {
	Type            common.ReporterType `mapstructure:"type"`
	WebhookEndpoint string              `mapstructure:"webhookEndpoint"`
	WebhookTimeout  time.Duration       `mapstructure:"webhookTimeout"`
	FileOutputPath  string              `mapstructure:"fileOutputPath"`
}
