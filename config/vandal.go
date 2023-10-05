package config

import "time"

type VandalConfig struct {
	Endpoint   string        `mapstructure:"endpoint"`
	CfgTimeout time.Duration `mapstructure:"cfgTimeout"`
}
