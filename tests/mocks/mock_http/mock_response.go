package mock_http

import "net/http"

// MockResponseWriter mock for http.ResponseWriter interface
type MockResponseWriter struct {
	HeaderMock      func() http.Header
	WriteMock       func([]byte) (int, error)
	WriteHeaderMock func(int)
	Headers         map[string][]string
	Body            []byte
	Status    		int
}

// Header mock for Header
func (m MockResponseWriter) Header() http.Header {
	if m.HeaderMock != nil {
		m.HeaderMock()
	}
	return m.Headers
}

// Write mock for Write
func (m *MockResponseWriter) Write(b []byte) (int, error) {
	if m.HeaderMock != nil {
		m.WriteMock(b)
	}
	m.Body = b
	return 0, nil
}

// WriteHeader mock for WriteHeader
func (m *MockResponseWriter) WriteHeader(i int) {
	if m.HeaderMock != nil {
		m.WriteHeaderMock(i)
	}
	m.Status = i
}
