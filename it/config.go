package it

import "github.com/dogefuzz/dogefuzz/config"

var GETH_CONFIG = config.GethConfig{
	NodeAddress: "http://localhost:51774",
	ChainID:     1900,
}

const SOLC_FOLDER = "/tmp/solc"
