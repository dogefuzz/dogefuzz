package oracle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DelegateOracleTestSuite struct {
	suite.Suite
}

// test detection of delegate call weakness with correct case
func (suite *DelegateOracleTestSuite) TestDetectCorrectCase() {
	snapshot := EventsSnapshot{
		Delegate: true,
	}

	oracle := DelegateOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "delegate call didn't detect weakness")
}

// test detection of delegate call weakness with wrong case
func (suite *DelegateOracleTestSuite) TestDetectWrongCase() {
	snapshot := EventsSnapshot{
		Delegate: false,
	}

	oracle := DelegateOracle{}
	assert.False(suite.T(), oracle.Detect(snapshot), "delegate call incorrectly detected weakness")
}

func TestDelegateOracleTestSuite(t *testing.T) {
	suite.Run(t, new(DelegateOracleTestSuite))
}
