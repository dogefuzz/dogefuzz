package oracle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TimestampDependencyOracleTestSuite struct {
	suite.Suite
}

// test detection of exception disorder call weakness with correct case 1
func (suite *TimestampDependencyOracleTestSuite) TestDetectCorrectCase1() {
	snapshot := EventsSnapshot{
		Timestamp:      true,
		StorageChanged: true,
	}

	oracle := TimestampDependencyOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with correct case 2
func (suite *TimestampDependencyOracleTestSuite) TestDetectCorrectCase2() {
	snapshot := EventsSnapshot{
		Timestamp:     true,
		EtherTransfer: true,
	}

	oracle := TimestampDependencyOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with correct case 3
func (suite *TimestampDependencyOracleTestSuite) TestDetectCorrectCase3() {
	snapshot := EventsSnapshot{
		Timestamp: true,
		SendOp:    true,
	}

	oracle := TimestampDependencyOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with wrong case 1
func (suite *TimestampDependencyOracleTestSuite) TestDetectWrongCase1() {
	snapshot := EventsSnapshot{
		Timestamp:      false,
		StorageChanged: true,
		EtherTransfer:  true,
		SendOp:         true,
	}

	oracle := TimestampDependencyOracle{}
	assert.False(suite.T(), oracle.Detect(snapshot), "exception disorder call incorrectly detected weakness")
}

// test detection of exception disorder call weakness with wrong case 2
func (suite *TimestampDependencyOracleTestSuite) TestDetectWrongCase2() {
	snapshot := EventsSnapshot{
		Timestamp:      true,
		StorageChanged: false,
		EtherTransfer:  false,
		SendOp:         false,
	}

	oracle := TimestampDependencyOracle{}
	assert.False(suite.T(), oracle.Detect(snapshot), "exception disorder call incorrectly detected weakness")
}

func TestTimestampDependencyOracleTestSuite(t *testing.T) {
	suite.Run(t, new(TimestampDependencyOracleTestSuite))
}
