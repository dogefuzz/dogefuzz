package solc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SolidityCompilerTestSuite struct {
	suite.Suite
}

func TestSolidityCompilerTestSuite(t *testing.T) {
	suite.Run(t, new(SolidityCompilerTestSuite))
}

func (s *SolidityCompilerTestSuite) TestCompileSource() {
	compiler := NewSolidityCompiler("/tmp/dogefuzz")

	contract1, err := compiler.CompileSource("HelloWorld", CONTRACT1)
	if err != nil {
		s.Fail(fmt.Sprintf("And error occurred: %s", err))
	}

	// const expectedDeploymentBytecode1 = "0x60c0604052600c60809081526b48656c6c6f20576f726c642160a01b60a05260009061002b90826100dd565b5034801561003857600080fd5b5061019c565b634e487b7160e01b600052604160045260246000fd5b600181811c9082168061006857607f821691505b60208210810361008857634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156100d857600081815260208120601f850160051c810160208610156100b55750805b601f850160051c820191505b818110156100d4578281556001016100c1565b5050505b505050565b81516001600160401b038111156100f6576100f661003e565b61010a816101048454610054565b8461008e565b602080601f83116001811461013f57600084156101275750858301515b600019600386901b1c1916600185901b1785556100d4565b600085815260208120601f198616915b8281101561016e5788860151825594840194600190910190840161014f565b508582101561018c5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61019a806101ab6000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063cfae321714610030575b600080fd5b61003861004e565b60405161004591906100dc565b60405180910390f35b6000805461005b9061012a565b80601f01602080910402602001604051908101604052809291908181526020018280546100879061012a565b80156100d45780601f106100a9576101008083540402835291602001916100d4565b820191906000526020600020905b8154815290600101906020018083116100b757829003601f168201915b505050505081565b600060208083528351808285015260005b81811015610109578581018301518582016040015282016100ed565b506000604082860101526040601f19601f8301168501019250505092915050565b600181811c9082168061013e57607f821691505b60208210810361015e57634e487b7160e01b600052602260045260246000fd5b5091905056fea26469706673582212209587fdcc9104e85579f4ee9820ef9720bd0f6a91225253a204a5a663e890075064736f6c63430008130033"
	// assert.Equal(s.T(), expectedDeploymentBytecode1, contract1.DeploymentBytecode)
	assert.Greater(s.T(), len(contract1.DeploymentBytecode), 2)
	const expectedAbiDefinition1 = "[{\"inputs\":[],\"name\":\"greet\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	assert.Equal(s.T(), expectedAbiDefinition1, contract1.AbiDefinition)

	contract2, err := compiler.CompileSource("HashForEther", CONTRACT2)
	if err != nil {
		s.Fail(fmt.Sprintf("And error occurred: %s", err))
	}
	// const expectedDeploymentBytecode2 = "0x608060405234801561001057600080fd5b5060e68061001f6000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166383ac4ae18114604d578063cc42e83a146061575b600080fd5b348015605857600080fd5b50605f6073565b005b348015606c57600080fd5b50605f60a2565b6040513390303180156108fc02916000818181858888f19350505050158015609f573d6000803e3d6000fd5b50565b63ffffffff33161560b257600080fd5b60b86073565b5600a165627a7a72305820a7e90b055e105f828a1ca91ba5438fdaad26e17ded294e9e8a3b12ac41a0d4f20029"
	assert.Greater(s.T(), len(contract2.DeploymentBytecode), 2)
	const expectedAbiDefinition2 = "[{\"constant\":false,\"inputs\":[],\"name\":\"_sendWinnings\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawWinnings\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	assert.Equal(s.T(), expectedAbiDefinition2, contract2.AbiDefinition)
}

func (s *SolidityCompilerTestSuite) TestDecompile_WithDefaultContracts() {
	currentDirectory, err := os.Getwd()
	assert.Nil(s.T(), err)

	contractFolder := filepath.Join(currentDirectory, "../../assets/contracts")
	solidityFiles, _ := os.ReadDir(contractFolder)

	for _, file := range solidityFiles {
		f, _ := os.ReadFile(filepath.Join(contractFolder, file.Name()))
		compiler := NewSolidityCompiler("/tmp/dogefuzz/")
		fileWithoutExtension := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
		contract, err := compiler.CompileSource(fileWithoutExtension, string(f))
		assert.Nil(s.T(), err, "error on "+file.Name())
		assert.NotEqual(s.T(), "0x", contract.DeploymentBytecode)
	}
}

const CONTRACT1 = `
pragma solidity ^0.8.13;

contract HelloWorld {
    string public greet = "Hello World!";
}
`

const CONTRACT2 = `
pragma solidity ^0.4.24;

contract HashForEther {

    function withdrawWinnings() {
        // Winner if the last 8 hex characters of the address are 0.
        require(uint32(msg.sender) == 0);
        _sendWinnings();
     }

     function _sendWinnings() {
         msg.sender.transfer(this.balance);
     }
}

`
