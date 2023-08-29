package config

type GethConfig struct {
	NodeAddress           string `mapstructure:"nodeAddress"`
	ChainID               int64  `mapstructure:"chainId"`
	MinGasLimit           uint64 `mapstructure:"minGasLimit"`
	DeployerAddress       string `mapstructure:"deployerAddress"`
	DeployerPrivateKeyHex string `mapstructure:"deployerPrivateKeyHex"`
	AgentAddress          string `mapstructure:"agentAddress"`
	AgentPrivateKeyHex    string `mapstructure:"agentPrivateKeyHex"`
}
