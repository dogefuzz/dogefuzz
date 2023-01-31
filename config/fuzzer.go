package config

import "github.com/dogefuzz/dogefuzz/pkg/common"

type FuzzerConfig struct {
	CritialInstructions []string     `mapstructure:"criticalInstructions"`
	BatchSize           int          `mapstructure:"batchSize"`
	SeedsSize           int          `mapstructure:"seedsSize"`
	Seeds               common.Seeds `mapstructure:"seeds"`
}
