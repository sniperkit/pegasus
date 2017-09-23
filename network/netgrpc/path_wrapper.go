package netgrpc

import "bitbucket.org/code_horse/pegasus/network"

// PathWrapper describes the router field. Each Field container a handler and a middleware which are function.
type PathWrapper struct {
	Handler    network.Handler
	Middleware network.Middleware
}
