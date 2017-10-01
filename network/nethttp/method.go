package nethttp

import (
	"net/http"
)

// Method describe the HTTP methods
type Method string

const (
	// Get Http Method
	Get Method = http.MethodGet
	// Head Http Method
	Head Method = http.MethodHead
	// Post Http Method
	Post Method = http.MethodPost
	// Put Http Method
	Put Method = http.MethodPut
	// Patch Http Method
	Patch Method = http.MethodPatch
	// Delete Http Method
	Delete Method = http.MethodDelete
	// Connect Http Method
	Connect Method = http.MethodConnect
	// Options Http Method
	Options Method = http.MethodOptions
	// Trace Http Method
	Trace Method = http.MethodTrace
)

// String return the string that is related with the specific method
func (m Method) String() string {
	return string(m)
}
