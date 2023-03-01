package config

type GethConfig struct {
	NodeAddress           string `mapstructure:"nodeAddress"`
	ChainID               int64  `mapstructure:"chainId"`
	DeployerAddress       string `mapstructure:"deployerAddress"`
	DeployerPrivateKeyHex string `mapstructure:"deployerPrivateKeyHex"`
	AgentAddress          string `mapstructure:"agentAddress"`
	AgentPrivateKeyHex    string `mapstructure:"agentPrivateKeyHex"`
}
