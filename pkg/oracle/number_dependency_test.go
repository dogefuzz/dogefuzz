package oracle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NumberDependencyOracleTestSuite struct {
	suite.Suite
}

// test detection of exception disorder call weakness with correct case 1
func (suite *NumberDependencyOracleTestSuite) TestDetectCorrectCase1() {
	snapshot := EventsSnapshot{
		BlockNumber:    true,
		StorageChanged: true,
	}

	oracle := NumberDependencyOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with correct case 2
func (suite *NumberDependencyOracleTestSuite) TestDetectCorrectCase2() {
	snapshot := EventsSnapshot{
		BlockNumber:   true,
		EtherTransfer: true,
	}

	oracle := NumberDependencyOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with correct case 3
func (suite *NumberDependencyOracleTestSuite) TestDetectCorrectCase3() {
	snapshot := EventsSnapshot{
		BlockNumber: true,
		SendOp:      true,
	}

	oracle := NumberDependencyOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with wrong case 1
func (suite *NumberDependencyOracleTestSuite) TestDetectWrongCase1() {
	snapshot := EventsSnapshot{
		BlockNumber:    false,
		StorageChanged: true,
		EtherTransfer:  true,
		SendOp:         true,
	}

	oracle := NumberDependencyOracle{}
	assert.False(suite.T(), oracle.Detect(snapshot), "exception disorder call incorrectly detected weakness")
}

// test detection of exception disorder call weakness with wrong case 2
func (suite *NumberDependencyOracleTestSuite) TestDetectWrongCase2() {
	snapshot := EventsSnapshot{
		BlockNumber:    true,
		StorageChanged: false,
		EtherTransfer:  false,
		SendOp:         false,
	}

	oracle := NumberDependencyOracle{}
	assert.False(suite.T(), oracle.Detect(snapshot), "exception disorder call incorrectly detected weakness")
}

func TestNumberDependencyOracleTestSuite(t *testing.T) {
	suite.Run(t, new(NumberDependencyOracleTestSuite))
}
