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
	contract, err := compiler.CompileSource("HelloWorld", VALID_SOLIDITY_FILE_WITH_NO_CONSTRUCTOR)
	assert.Nil(s.T(), err)

	address, tx, err := deployer.Deploy(context.Background(), contract)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), address)
	assert.NotEmpty(s.T(), tx)

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
	contract, err := compiler.CompileSource("HelloWorld", VALID_SOLIDITY_FILE_WITH_CONSTRUCTOR)
	assert.Nil(s.T(), err)

	arg := gofakeit.Word()
	address, tx, err := deployer.Deploy(context.Background(), contract, arg)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), address)
	assert.NotEmpty(s.T(), tx)

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

// func (s *GethDeployerIntegrationTestSuite) TestDeploy_WithBenchmarkContracts() {

// 	folder := "/home/imedeiros/workspace/dogefuzz/dogefuzz/test/resources/contracts"
// 	solidityFiles, _ := ioutil.ReadDir(folder)

// 	deployer, err := NewDeployer(it.GETH_CONFIG)
// 	assert.Nil(s.T(), err)

// 	for _, file := range solidityFiles {

// 		f, _ := ioutil.ReadFile("/home/imedeiros/workspace/dogefuzz/dogefuzz/test/resources/contracts/" + file.Name())
// 		compiler := solc.NewSolidityCompiler("/tmp/dogefuzz/")
// 		fileWithoutExtension := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
// 		contract, err := compiler.CompileSource(fileWithoutExtension, string(f))
// 		assert.Nil(s.T(), err, fmt.Sprintf("error on %s", file.Name()))
// 		assert.NotEqual(s.T(), "0x", contract.CompiledCode)

// 		parsedAbi, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
// 		assert.Nil(s.T(), err, fmt.Sprintf("error on %s", file.Name()))

// 		inputs := make([]interface{}, 0)
// 		for _, parameter := range parsedAbi.Constructor.Inputs {
// 			handler, err := solidity.GetTypeHandler(parameter.Type)
// 			assert.Nil(s.T(), err, fmt.Sprintf("error on contract %s", file.Name()))

// 			handler.Generate()
// 			inputs = append(inputs, handler.GetValue())
// 		}
// 		address, tx, err := deployer.Deploy(context.Background(), contract, inputs...)
// 		assert.Nil(s.T(), err, fmt.Sprintf("error on %s", file.Name()))
// 		assert.NotEmpty(s.T(), address)
// 		assert.NotEmpty(s.T(), tx)
// 	}
// }

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
