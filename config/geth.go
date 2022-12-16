package config

type GethConfig struct {
	NodeAddress           string `mapstructure:"GETH_NODE_ADDRESS"`
	ChainID               int64  `mapstructure:"GETH_CHAIN_ID"`
	DeployerPrivateKeyHex string `mapstructure:"GETH_DEPLOYER_PRIVATE_KEY"`
	AgentPrivateKeyHex    string `mapstructure:"GETH_AGENT_PRIVATE_KEY"`
}
