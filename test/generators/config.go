package generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/config"
)

func ConfigGen() *config.Config {
	return &config.Config{
		DatabaseName: gofakeit.LetterN(255),
		GethConfig:   GethConfigGen(),
	}
}

func GethConfigGen() config.GethConfig {
	return config.GethConfig{
		NodeAddress:           gofakeit.URL(),
		ChainID:               gofakeit.Int64(),
		DeployerPrivateKeyHex: gofakeit.HexUint32(),
		AgentPrivateKeyHex:    gofakeit.HexUint32(),
	}
}
