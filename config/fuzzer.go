package config

import "github.com/dogefuzz/dogefuzz/pkg/solidity"

type FuzzerConfig struct {
	CritialInstructions []string                             `mapstructure:"criticalInstructions"`
	BatchSize           int                                  `mapstructure:"batchSize"`
	SeedsSize           int                                  `mapstructure:"seedsSize"`
	Seeds               map[solidity.TypeIdentifier][]string `mapstructure:"seeds"`
}
