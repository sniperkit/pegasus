package mock_http

import "net/http"

// MockHTTPClient is a mock for http.Client struct
type MockHTTPClient struct {
	DoMock func(req *http.Request) (*http.Response, error)
}

// Do mock for http.Client Do
func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.DoMock != nil {
		return c.DoMock(req)
	}
	return &http.Response{}, nil
}
