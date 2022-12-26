package config

import "github.com/dogefuzz/dogefuzz/pkg/solidity"

type FuzzerConfig struct {
	CritialInstructions []string                             `json:"critialInstructions"`
	BatchSize           int                                  `json:"batchSize"`
	SeedsSize           int                                  `json:"seedsSize"`
	Seeds               map[solidity.TypeIdentifier][]string `json:"seeds"`
}
