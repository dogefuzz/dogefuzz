package oracle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExceptionDisorderOracleTestSuite struct {
	suite.Suite
}

// test detection of exception disorder call weakness with correct case
func (suite *ExceptionDisorderOracleTestSuite) TestDetectCorrectCase() {
	snapshot := EventsSnapshot{
		ExceptionDisorder: true,
	}

	oracle := ExceptionDisorderOracle{}
	assert.True(suite.T(), oracle.Detect(snapshot), "exception disorder call didn't detect weakness")
}

// test detection of exception disorder call weakness with wrong case
func (suite *ExceptionDisorderOracleTestSuite) TestDetectWrongCase() {
	snapshot := EventsSnapshot{
		ExceptionDisorder: false,
	}

	oracle := ExceptionDisorderOracle{}
	assert.False(suite.T(), oracle.Detect(snapshot), "exception disorder call incorrectly detected weakness")
}

func TestExceptionDisorderOracleTestSuite(t *testing.T) {
	suite.Run(t, new(ExceptionDisorderOracleTestSuite))
}
