package oracle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GaslessSendOracleTestSuite struct {
	suite.Suite
}

// test detection of exception disorder call weakness with correct case
func (suite *GaslessSendOracleTestSuite) TestDetectCorrectCase() {
	snapshot := EventsSnapshot{
		GaslessSend: true,
	}

	oracle := GaslessSendOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with wrong case
func (suite *GaslessSendOracleTestSuite) TestDetectWrongCase() {
	snapshot := EventsSnapshot{
		GaslessSend: false,
	}

	oracle := GaslessSendOracle{}
	assert.False(suite.T(), oracle.Detect(snapshot), "exception disorder call incorrectly detected weakness")
}

func TestGaslessSendOracleTestSuite(t *testing.T) {
	suite.Run(t, new(GaslessSendOracleTestSuite))
}
