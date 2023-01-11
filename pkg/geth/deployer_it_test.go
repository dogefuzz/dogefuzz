package geth

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/dogefuzz/dogefuzz/test/it"
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

func (s *GethDeployerIntegrationTestSuite) TestDeploy_ShouldDeployContractInGethNode_WhenProvidedAValidContractWithNoConstructor() {
	deployer, err := NewDeployer(it.GETH_CONFIG)
	assert.Nil(s.T(), err)

	compiler := solc.NewSolidityCompiler(it.SOLC_FOLDER)
	contract, err := compiler.CompileSource(VALID_SOLIDITY_FILE_WITH_NO_CONSTRUCTOR)
	assert.Nil(s.T(), err)

	address, err := deployer.Deploy(context.Background(), contract)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), address)

	client, err := ethclient.Dial(it.GETH_CONFIG.NodeAddress)
	assert.Nil(s.T(), err)

	parsedAbi, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	assert.Nil(s.T(), err)

	bindedContract := bind.NewBoundContract(common.HexToAddress(address), parsedAbi, client, client, client)
	var results []interface{}
	err = bindedContract.Call(nil, &results, "greet")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "Hello World!", results[0].(string))
}

func (s *GethDeployerIntegrationTestSuite) TestDeploy_ShouldDeployContractInGethNode_WhenProvidedAValidContractWithConstructor() {
	deployer, err := NewDeployer(it.GETH_CONFIG)
	assert.Nil(s.T(), err)

	compiler := solc.NewSolidityCompiler(it.SOLC_FOLDER)
	contract, err := compiler.CompileSource(VALID_SOLIDITY_FILE_WITH_CONSTRUCTOR)
	assert.Nil(s.T(), err)

	arg := gofakeit.Word()
	address, err := deployer.Deploy(context.Background(), contract, arg)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), address)

	client, err := ethclient.Dial(it.GETH_CONFIG.NodeAddress)
	assert.Nil(s.T(), err)

	parsedAbi, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	assert.Nil(s.T(), err)

	bindedContract := bind.NewBoundContract(common.HexToAddress(address), parsedAbi, client, client, client)
	var results []interface{}
	err = bindedContract.Call(nil, &results, "greet")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("Hello, %s!", arg), results[0].(string))
}

const VALID_SOLIDITY_FILE_WITH_NO_CONSTRUCTOR = `
// SPDX-License-Identifier: MIT
// compiler version must be greater than or equal to 0.8.13 and less than 0.9.0
pragma solidity ^0.4.26;

contract HelloWorld {
    string public greet = "Hello World!";
}
`

const VALID_SOLIDITY_FILE_WITH_CONSTRUCTOR = `
// SPDX-License-Identifier: MIT
// compiler version must be greater than or equal to 0.8.13 and less than 0.9.0
pragma solidity ^0.4.26;

contract HelloWorld {

    string private _name;

    constructor(string name) public {
        _name = name;
    }

    function greet() public view returns (string memory) {
        return string(abi.encodePacked("Hello, ", _name, "!"));
    }
}
`
