package nethttp

import (
	"net/http"
	"io/ioutil"
)

// NewRequest http.NewRequest
var NewRequest = http.NewRequest

// ReadAll reads the byte
var ReadAll = ioutil.ReadAll
