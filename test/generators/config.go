package generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func ConfigGen() *config.Config {
	return &config.Config{
		DatabaseName: gofakeit.LetterN(255),
		GethConfig:   GethConfigGen(),
		FuzzerConfig: FuzzerConfigGen(),
	}
}

func FuzzerConfigGen() config.FuzzerConfig {
	criticalInstructionsCount := gofakeit.Number(1, 10)
	criticalInstructions := make([]string, criticalInstructionsCount)
	for idx := 0; idx < int(criticalInstructionsCount); idx++ {
		criticalInstructions[idx] = gofakeit.LetterN(255)
	}

	seeds := make(map[common.TypeIdentifier][]string)
	return config.FuzzerConfig{
		CritialInstructions:          criticalInstructions,
		BatchSize:                    int(gofakeit.Int32()),
		SeedsSize:                    int(gofakeit.Int32()),
		Seeds:                        seeds,
		TransactionTimeout:           gofakeit.Date().Sub(gofakeit.Date()),
		PendingTransactionsThreshold: int(gofakeit.Int32()),
	}
}

func GethConfigGen() config.GethConfig {
	return config.GethConfig{
		NodeAddress:           gofakeit.URL(),
		ChainID:               gofakeit.Int64(),
		MinGasLimit:           gofakeit.Uint64(),
		DeployerPrivateKeyHex: gofakeit.HexUint32(),
		AgentPrivateKeyHex:    gofakeit.HexUint32(),
	}
}
