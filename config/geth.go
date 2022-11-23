package config

type GethConfig struct {
	NodeAddress string `mapstructure:"GETH_NODE_ADDRESS"`
	ChainID     int64  `mapstructure:"GETH_CHAIN_ID"`
}
