package trangrpc

import "bitbucket.org/code_horse/pegasus/network"

type PathWrapper struct {
	Handler    network.Handler
	Middleware network.Middleware
}
