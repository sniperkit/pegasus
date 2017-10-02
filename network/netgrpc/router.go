package netgrpc

import "bitbucket.org/code_horse/pegasus/network"

// Router struct describes the paths, handlers and the middleware. This is a mapper in order
// server to know when to call each handler.
type Router struct {

	// PathWrapper is the mapper which has a key that points to a unique path (url) and is referring to the handler and
	// the middleware
	PathsWrapper map[string]*PathWrapper
}

// NewRouter is a construct for Router struct which initialize and returns a *Router object.
func NewRouter() *Router {
	return &Router{
		PathsWrapper: make(map[string]*PathWrapper),
	}
}

// Add is the function which add a new PathWrapper in Router struct. It gets the path as a unique key.
func (r *Router) Add(path string, handler network.Handler, middleware network.Middleware) {
	r.PathsWrapper[path] = &PathWrapper{
		Handler:    handler,
		Middleware: middleware,
	}
}
