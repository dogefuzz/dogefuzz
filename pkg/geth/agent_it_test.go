package geth

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/dogefuzz/dogefuzz/test/it"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

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

	deployer, err := NewDeployer(zap.NewNop(), it.GETH_CONFIG)
	assert.Nil(s.T(), err)

	compiler := solc.NewSolidityCompiler(it.SOLC_FOLDER)
	contract, err := compiler.CompileSource("HelloWorld", solidityFile)
	assert.Nil(s.T(), err)

	address, tx, err := deployer.Deploy(context.Background(), contract, gofakeit.Word())
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), address)
	assert.NotEmpty(s.T(), tx)

	agent, err := NewAgent(it.GETH_CONFIG)
	assert.Nil(s.T(), err)

	privateKey := common.RandomChoice([]string{it.GETH_CONFIG.AgentPrivateKeyHex, it.GETH_CONFIG.DeployerPrivateKeyHex})
	wallet, err := NewWalletFromPrivateKeyHex(privateKey)
	assert.Nil(s.T(), err)

	newWord := gofakeit.Word()
	value := big.NewInt(0)
	tx, err = agent.Send(context.Background(), wallet, contract, "setName", value, newWord)
	assert.Nil(s.T(), err)

	client, err := ethclient.Dial(it.GETH_CONFIG.NodeAddress)
	assert.Nil(s.T(), err)

	var receipt *types.Receipt
	for {
		receipt, err = client.TransactionReceipt(context.Background(), gethcommon.HexToHash(tx))
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

	bindedContract := bind.NewBoundContract(gethcommon.HexToAddress(address), parsedAbi, client, client, client)
	var results []interface{}
	err = bindedContract.Call(nil, &results, "greet")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("Hello, %s!", newWord), results[0].(string))
}

// func (s *GethAgentIntegrationTestSuite) TestSend_WithBenchmarkContracts() {

// 	folder := "/home/imedeiros/workspace/dogefuzz/dogefuzz/test/resources/contracts"
// 	solidityFiles, _ := ioutil.ReadDir(folder)

// 	deployer, err := NewDeployer(it.GETH_CONFIG)
// 	assert.Nil(s.T(), err)

// 	for _, file := range solidityFiles {

// 		f, _ := ioutil.ReadFile("/home/imedeiros/workspace/dogefuzz/dogefuzz/test/resources/contracts/" + file.Name())
// 		compiler := solc.NewSolidityCompiler("/tmp/dogefuzz/")
// 		fileWithoutExtension := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
// 		contract, err := compiler.CompileSource(fileWithoutExtension, string(f))
// 		assert.Nil(s.T(), err, fmt.Sprintf("error on contract %s", file.Name()))
// 		assert.NotEqual(s.T(), "0x", contract.DeploymentBytecode)

// 		parsedAbi, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
// 		assert.Nil(s.T(), err, fmt.Sprintf("error on contract %s", file.Name()))

// 		inputs := make([]interface{}, 0)
// 		for _, parameter := range parsedAbi.Constructor.Inputs {
// 			handler, err := solidity.GetTypeHandler(parameter.Type)
// 			assert.Nil(s.T(), err, fmt.Sprintf("error on contract %s", file.Name()))

// 			handler.Generate()
// 			inputs = append(inputs, handler.GetValue())
// 		}
// 		address, err := deployer.Deploy(context.Background(), contract, inputs...)
// 		assert.Nil(s.T(), err, fmt.Sprintf("error on contract %s", file.Name()))
// 		assert.NotEmpty(s.T(), address)

// 		agent, err := NewAgent(it.GETH_CONFIG)
// 		assert.Nil(s.T(), err, fmt.Sprintf("error on contract %s", file.Name()))

// 		for _, method := range parsedAbi.Methods {
// 			if !method.Payable {
// 				continue
// 			}
// 			inputs = make([]interface{}, 0)
// 			for _, parameter := range method.Inputs {
// 				handler, err := solidity.GetTypeHandler(parameter.Type)
// 				assert.Nil(s.T(), err, fmt.Sprintf("error on contract %s", file.Name()))

// 				handler.Generate()
// 				inputs = append(inputs, handler.GetValue())
// 			}

// 			nonce, err := agent.GetNonce(context.Background())
// 			assert.Nil(s.T(), err)

// 			tx, err := agent.Send(context.Background(), nonce, contract, method.Name, inputs...)
// 			assert.NotEmpty(s.T(), tx, fmt.Sprintf("error on contract %s in method %s", file.Name(), method.Name))
// 			assert.Nil(s.T(), err, fmt.Sprintf("error on contract %s in method %s", file.Name(), method.Name))

// 			client, err := ethclient.Dial(it.GETH_CONFIG.NodeAddress)
// 			assert.Nil(s.T(), err, fmt.Sprintf("error on contract %s in method %s", file.Name(), method.Name))

// 			if tx != "" {
// 				var receipt *types.Receipt
// 				for {
// 					receipt, err = client.TransactionReceipt(context.Background(), common.HexToHash(tx))
// 					if err != nil {
// 						if err != ethereum.NotFound {
// 							assert.Nil(s.T(), err, fmt.Sprintf("error on contract %s in method %s", file.Name(), method.Name))
// 						}
// 					} else {
// 						break
// 					}

// 					time.Sleep(1 * time.Second)
// 				}
// 				assert.NotNil(s.T(), receipt)
// 			} else {
// 				break
// 			}
// 		}
// 	}
// }
