package grpc_transport

import "bitbucket.org/code_horse/pegasus/network"

type Router struct {
	PathsWrapper map[string]*PathWrapper
}

func (r *Router) Add(path string, handler network.Handler, middleware network.Middleware) {
	r.PathsWrapper[path] = &PathWrapper{
		Handler: handler,
		Middleware: middleware,
	}
}
