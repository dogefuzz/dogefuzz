package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HTTPClientMock struct {
	mock.Mock
}

func (m *HTTPClientMock) Do(r *http.Request) (*http.Response, error) {
	args := m.Called(r)
	return args.Get(0).(*http.Response), args.Error(1)
}
