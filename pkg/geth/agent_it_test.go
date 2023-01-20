package geth

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/dogefuzz/dogefuzz/test/it"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GethAgentIntegrationTestSuite struct {
	suite.Suite
}

func TestGethAgentIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(GethAgentIntegrationTestSuite))
}

func (s *GethAgentIntegrationTestSuite) TestDeploy_ShouldDeployContractInGethNode_WhenProvidedAValidContractWithNoConstructor() {
	const solidityFile = `
// SPDX-License-Identifier: MIT
// compiler version must be greater than or equal to 0.8.13 and less than 0.9.0
pragma solidity ^0.4.26;

contract HelloWorld {
	string private name;

	constructor (string _name) public {
		name = _name;
	}

	function greet() public view returns (string memory) {
		return string(abi.encodePacked("Hello, ", name, "!"));
	}

	function setName(string _name) public payable {
		name = _name;
	}
}
`

	deployer, err := NewDeployer(it.GETH_CONFIG)
	assert.Nil(s.T(), err)

	compiler := solc.NewSolidityCompiler(it.SOLC_FOLDER)
	contract, err := compiler.CompileSource("HelloWorld", solidityFile)
	assert.Nil(s.T(), err)

	address, err := deployer.Deploy(context.Background(), contract, gofakeit.Word())
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), address)

	agent, err := NewAgent(it.GETH_CONFIG)
	assert.Nil(s.T(), err)

	newWord := gofakeit.Word()
	tx, err := agent.Send(context.Background(), contract, "setName", newWord)
	assert.Nil(s.T(), err)

	client, err := ethclient.Dial(it.GETH_CONFIG.NodeAddress)
	assert.Nil(s.T(), err)

	var receipt *types.Receipt
	for {
		receipt, err = client.TransactionReceipt(context.Background(), common.HexToHash(tx))
		if err != nil {
			if err != ethereum.NotFound {
				assert.Nil(s.T(), err)
			}
		} else {
			break
		}

		time.Sleep(1 * time.Second)
	}
	assert.NotNil(s.T(), receipt)

	parsedAbi, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	assert.Nil(s.T(), err)

	bindedContract := bind.NewBoundContract(common.HexToAddress(address), parsedAbi, client, client, client)
	var results []interface{}
	err = bindedContract.Call(nil, &results, "greet")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("Hello, %s!", newWord), results[0].(string))
}
