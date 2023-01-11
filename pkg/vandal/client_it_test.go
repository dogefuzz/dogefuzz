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
	contract, _ := compiler.CompileSource(VALID_SOLIDITY_FILE)

	c := NewVandalClient("http://localhost:51243")
	blocks, functions, err := c.Decompile(context.Background(), contract.CompiledCode)
	assert.Equal(s.T(), 30, len(blocks))
	assert.Equal(s.T(), 1, len(functions))
	assert.Nil(s.T(), err)
}
