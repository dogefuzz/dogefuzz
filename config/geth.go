package config

type GethConfig struct {
	NodeAddress           string `mapstructure:"nodeAddress"`
	ChainID               int64  `mapstructure:"chainId"`
	DeployerPrivateKeyHex string `mapstructure:"deployerPrivateKeyHex"`
	AgentPrivateKeyHex    string `mapstructure:"agentPrivateKeyHex"`
}
