package vandal

import (
	"context"
	"testing"

	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VandalClientIntegrationTestSuite struct {
	suite.Suite
}

func TestVandalClientIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(VandalClientIntegrationTestSuite))
}

func (s *VandalClientIntegrationTestSuite) TestDecompile_WhenPassingAValidSource_ReturnBlocksAndFunctions() {
	const VALID_SOLIDITY_FILE = `
// SPDX-License-Identifier: MIT
// compiler version must be greater than or equal to 0.8.13 and less than 0.9.0
pragma solidity ^0.4.18;

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
	compiler := solc.NewSolidityCompiler("/tmp/dogefuzz/")
	contract, _ := compiler.CompileSource("HelloWorld", VALID_SOLIDITY_FILE)

	c := NewVandalClient("http://localhost:51243")
	blocks, functions, err := c.Decompile(context.Background(), contract.CompiledCode)
	assert.Equal(s.T(), 30, len(blocks))
	assert.Equal(s.T(), 2, len(functions))
	assert.Nil(s.T(), err)
}

// func (s *VandalClientIntegrationTestSuite) TestDecompile_WithBenchmarkContracts() {

// 	folder := "/home/imedeiros/workspace/dogefuzz/dogefuzz/test/resources/contracts"
// 	solidityFiles, _ := ioutil.ReadDir(folder)
// 	for _, file := range solidityFiles {

// 		f, _ := ioutil.ReadFile("/home/imedeiros/workspace/dogefuzz/dogefuzz/test/resources/contracts/" + file.Name())
// 		compiler := solc.NewSolidityCompiler("/tmp/dogefuzz/")
// 		fileWithoutExtension := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
// 		contract, err := compiler.CompileSource(fileWithoutExtension, string(f))
// 		assert.Nil(s.T(), err, fmt.Sprintf("error while compiling the code for %s contract", file.Name()))

// 		c := NewVandalClient("http://localhost:51243")
// 		blocks, functions, err := c.Decompile(context.Background(), contract.CompiledCode)
// 		assert.Greater(s.T(), len(blocks), 0, fmt.Sprintf("the contract %s should have more than one block", fileWithoutExtension))
// 		assert.GreaterOrEqual(s.T(), len(functions), 0, fmt.Sprintf("the contract %s should have more than one function", fileWithoutExtension))
// 		assert.Nil(s.T(), err)
// 	}
// }
