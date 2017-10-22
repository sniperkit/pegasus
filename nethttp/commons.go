package nethttp

import (
	"io/ioutil"
	"net/http"
)

// NewRequest http.NewRequest
var NewRequest = http.NewRequest

// ReadAll reads the byte
var ReadAll = ioutil.ReadAll
