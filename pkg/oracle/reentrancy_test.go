package oracle

import (
	"testing"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReentrancyOracleTestSuite struct {
	suite.Suite
}

// test detection of exception disorder call weakness with correct case 1
func (suite *ReentrancyOracleTestSuite) TestDetectCorrectCase1() {
	snapshot := common.EventsSnapshot{
		Reentrancy:     true,
		StorageChanged: true,
	}

	oracle := ReentrancyOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with correct case 2
func (suite *ReentrancyOracleTestSuite) TestDetectCorrectCase2() {
	snapshot := common.EventsSnapshot{
		Reentrancy:    true,
		EtherTransfer: true,
	}

	oracle := ReentrancyOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with correct case 3
func (suite *ReentrancyOracleTestSuite) TestDetectCorrectCase3() {
	snapshot := common.EventsSnapshot{
		Reentrancy: true,
		SendOp:     true,
	}

	oracle := ReentrancyOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with wrong case 1
func (suite *ReentrancyOracleTestSuite) TestDetectWrongCase1() {
	snapshot := common.EventsSnapshot{
		Reentrancy:     false,
		StorageChanged: true,
		EtherTransfer:  true,
		SendOp:         true,
	}

	oracle := ReentrancyOracle{}
	assert.False(suite.T(), oracle.Detect(snapshot), "exception disorder call incorrectly detected weakness")
}

// test detection of exception disorder call weakness with wrong case 2
func (suite *ReentrancyOracleTestSuite) TestDetectWrongCase2() {
	snapshot := common.EventsSnapshot{
		Reentrancy:     true,
		StorageChanged: false,
		EtherTransfer:  false,
		SendOp:         false,
	}

	oracle := ReentrancyOracle{}
	assert.False(suite.T(), oracle.Detect(snapshot), "exception disorder call incorrectly detected weakness")
}

func TestReentrancyOracleTestSuite(t *testing.T) {
	suite.Run(t, new(ReentrancyOracleTestSuite))
}
