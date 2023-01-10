package config

import "github.com/dogefuzz/dogefuzz/pkg/reporter"

type ReporterConfig struct {
	Type             reporter.ReporterType `mapstructure:"type"`
	CallbackEndpoint string                `mapstructure:"callbackEndpoint"`
}
