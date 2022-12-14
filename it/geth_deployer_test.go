package it

import (
	"context"
	"strings"
	"testing"

	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GethDeployerIntegrationTestSuite struct {
	suite.Suite
}

func TestGethDeployerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(GethDeployerIntegrationTestSuite))
}

func (s *GethDeployerIntegrationTestSuite) TestDeploy_ShouldDeployContractInGethNode_WhenProvidedAValidContract() {
	wallet, err := geth.NewWalletFromPrivateKeyHex("f33ff13222d9141bcfe072f4c148026bf0187a3ca1f7c4063a7f3e4aff6591a5")
	assert.Nil(s.T(), err)

	deployer, err := geth.NewDeployer(GETH_CONFIG, wallet)
	assert.Nil(s.T(), err)

	compiler := solc.NewSolidityCompiler(SOLC_FOLDER)
	contract, err := compiler.CompileSource(VALID_SOLIDITY_FILE)
	assert.Nil(s.T(), err)

	address, err := deployer.Deploy(context.Background(), contract)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), address)

	client, err := ethclient.Dial(GETH_CONFIG.NodeAddress)
	assert.Nil(s.T(), err)

	parsedAbi, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	assert.Nil(s.T(), err)

	bindedContract := bind.NewBoundContract(common.HexToAddress(address), parsedAbi, client, client, client)
	var results []interface{}
	err = bindedContract.Call(nil, &results, "greet")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "Hello World!", results[0].(string))

}

const VALID_SOLIDITY_FILE = `
// SPDX-License-Identifier: MIT
// compiler version must be greater than or equal to 0.8.13 and less than 0.9.0
pragma solidity ^0.4.18;

contract HelloWorld {
    string public greet = "Hello World!";
}
`
