package config

type GethConfig struct {
	NodeAddress           string `json:"nodeAddress"`
	ChainID               int64  `json:"chainId"`
	DeployerPrivateKeyHex string `json:"deployerPrivateKeyHex"`
	AgentPrivateKeyHex    string `json:"agentPrivateKeyHex"`
}
