package reporter

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/dogefuzz/dogefuzz/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WebhookReporterTestSuite struct {
	suite.Suite

	httpClientMock *mocks.HTTPClientMock
}

func TestWebhookReporterTestSuite(t *testing.T) {
	suite.Run(t, new(WebhookReporterTestSuite))
}

func (s *WebhookReporterTestSuite) SetupTest() {
	s.httpClientMock = new(mocks.HTTPClientMock)
}

func (s *WebhookReporterTestSuite) TestSendOutput_ShouldReturnNil_WhenSendValidRequestAndReceiveAValidResponse() {
	url := gofakeit.URL()
	timeout := time.Duration(gofakeit.Number(30, 60)) * time.Minute
	report := generators.TaskReportGen()
	r := NewWebhookReporter(s.httpClientMock, url, timeout)
	expectedResponse := &http.Response{Body: io.NopCloser(strings.NewReader("")), StatusCode: 200}
	s.httpClientMock.On("Do", mock.AnythingOfType("*http.Request")).Return(expectedResponse, nil)

	err := r.SendOutput(context.Background(), report)

	s.httpClientMock.AssertExpectations(s.T())
	assert.Nil(s.T(), err)
}

func (s *WebhookReporterTestSuite) TestSendOutput_ShouldReturnErrNonSuccessResponse_WhenSendValidRequestAndReceiveANon200Response() {
	url := gofakeit.URL()
	timeout := time.Duration(gofakeit.Number(30, 60)) * time.Minute
	report := generators.TaskReportGen()
	r := NewWebhookReporter(s.httpClientMock, url, timeout)
	expectedResponse := &http.Response{Body: io.NopCloser(strings.NewReader("")), StatusCode: 500}
	s.httpClientMock.On("Do", mock.AnythingOfType("*http.Request")).Return(expectedResponse, nil)

	err := r.SendOutput(context.Background(), report)

	s.httpClientMock.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), ErrNonSuccessResponse, err)
}

func (s *WebhookReporterTestSuite) TestSendOutput_ShouldReturnError_WhenClientReturnsAnError() {
	url := gofakeit.URL()
	timeout := time.Duration(gofakeit.Number(30, 60)) * time.Minute
	report := generators.TaskReportGen()
	r := NewWebhookReporter(s.httpClientMock, url, timeout)
	expectedError := errors.New("mocked error")
	s.httpClientMock.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, expectedError)

	err := r.SendOutput(context.Background(), report)

	s.httpClientMock.AssertExpectations(s.T())
	assert.ErrorIs(s.T(), expectedError, err)
}
