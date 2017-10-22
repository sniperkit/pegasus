package netgrpc

import (
	"github.com/cpapidas/pegasus/peg"
)

// PathWrapper describes the router field. Each Field contains a handler and a middleware which are functions.
type PathWrapper struct {
	Handler    peg.Handler
	Middleware peg.Middleware
}
