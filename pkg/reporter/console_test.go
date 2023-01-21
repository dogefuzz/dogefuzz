package reporter

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConsoleReporterTestSuite struct {
	suite.Suite

	httpClientMock *mocks.HTTPClientMock
}

func TestConsoleReporterTestSuite(t *testing.T) {
	suite.Run(t, new(ConsoleReporterTestSuite))
}

func (s *ConsoleReporterTestSuite) SetupTest() {
	s.httpClientMock = new(mocks.HTTPClientMock)
}

func (s *ConsoleReporterTestSuite) TestSendOutput_ShouldReturnNil_WhenSendValidRequestAndReceiveAValidResponse() {
	report := generators.TaskReportGen()
	buffer := new(bytes.Buffer)
	buffer.WriteString("********** FUZZING EXECUTION RESULT **********\n")
	buffer.WriteString(fmt.Sprintf("Time Elapsed: %v\n", report.TimeElapsed))
	buffer.WriteString(fmt.Sprintf("Contract Name: %s\n", report.ContractName))
	buffer.WriteString(fmt.Sprintf("Coverage: %d\n", report.Coverage))
	buffer.WriteString(fmt.Sprintf("Min Distance: %d\n", report.MinDistance))
	buffer.WriteString(fmt.Sprintf("Transactions: %d\n", len(report.Transactions)))
	buffer.WriteString("Weakneses Found:\n")
	for _, weakness := range report.DetectedWeaknesses {
		buffer.WriteString(fmt.Sprintf("\t- %s\n", weakness))
	}
	expectedOutput := buffer.String()

	var output bytes.Buffer
	r := NewConsoleReporter(&output)

	r.SendOutput(context.Background(), report)

	assert.Equal(s.T(), expectedOutput, output.String())
}
