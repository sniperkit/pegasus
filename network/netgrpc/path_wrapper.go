package netgrpc

import "github.com/cpapidas/pegasus/network"

// PathWrapper describes the router field. Each Field contains a handler and a middleware which are functions.
type PathWrapper struct {
	Handler    network.Handler
	Middleware network.Middleware
}
