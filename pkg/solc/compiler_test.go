package solc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CompilerTestSuite struct {
	suite.Suite
}

func TestCompilerTestSuite(t *testing.T) {
	suite.Run(t, new(CompilerTestSuite))
}

func (s *CompilerTestSuite) TestSource() {
	// contracts1, err := CompileSource("", CONTRACT1)
	// if err != nil {
	// 	s.Fail(fmt.Sprintf("And error occurred: %s", err))
	// }

	contracts2, err := CompileSource("", CONTRACT2)
	if err != nil {
		s.Fail(fmt.Sprintf("And error occurred: %s", err))
	}

	// fmt.Println("CONTRACT_ONE")
	// for key, element := range contracts1 {
	// 	fmt.Printf("%s => %s", key, element.RuntimeCode)
	// }
	fmt.Println("CONTRACT_TWO")
	for key, element := range contracts2 {
		fmt.Printf("%s => %s", key, element.Info.Source)
	}
	assert.True(s.T(), true)
}

func (s *CompilerTestSuite) TestExtractVersion() {
	version1 := ExtractVersion(CONTRACT1)
	version2 := ExtractVersion(CONTRACT2)

	assert.Equal(s.T(), "0.8.13", version1)
	assert.Equal(s.T(), "0.4.24", version2)
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
