package it

import "github.com/dogefuzz/dogefuzz/config"

var GETH_CONFIG = config.GethConfig{
	NodeAddress:           "http://localhost:51774",
	ChainID:               1900,
	DeployerPrivateKeyHex: "1c8d8e900c1b8c6554d995e172c3f58ebaf0e035be4f597e89aa3599cd970d9b",
	AgentPrivateKeyHex:    "39a7089e3f7e093b900bc8e98e6e5cba4639cd04d3e944ce1ddd2ca1595b7b87",
}

const SOLC_FOLDER = "/tmp/solc"
