package mockgrpc

import (
	"bitbucket.org/code_horse/pegasus/network"
	"bitbucket.org/code_horse/pegasus/network/netgrpc"
)

// Server is an empty struct that helps us test the real GRPC Server struct
type Server struct {
	router *netgrpc.Router
}

// NewServer returns a new mock *netgrpc.Server
func NewServer(router *netgrpc.Router) netgrpc.IServer {
	return &Server{router: router}
}

// Serve function checks throw a panic if address is nil.
func (s *Server) Serve(address string) {

	if address == "" {
		panic("Address should not be nil")
	}
}

// Listen throw a panic if params is not stetted
func (s *Server) Listen(path string, handler network.Handler, middleware network.Middleware) {

	if path == "" || handler == nil || middleware == nil {
		panic("params should not be nil")
	}

}

// ListenTo throw a panic if params is not stetted
func (s *Server) ListenTo(properties *netgrpc.Properties, handler network.Handler, middleware network.Middleware) {
	if handler == nil || middleware == nil {
		panic("params should not be nil")
	}
}

// Stop Mocked in order to implement the netgrpc.IServer interface
func (s *Server) Stop() {}
