package tranhttp

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Server struct is responsible for http server. It manages connections and configuration may needed in order to ensure
// that the http server works properly
type Server struct {

	// Router is responsible for handler and middleware
	Router *mux.Router
}

// NewServer is constructor of Server struct. It initialize and return a Server object. It get a *mux.Router as
// parameter which could not be nil
func NewServer(router *mux.Router) *Server {

	if router == nil {
		router = mux.NewRouter()
	}

	return &Server{Router: router}
}

// Serve function start the server for a specific part and port
func (s *Server) Serve(path string) {
	err := http.ListenAndServe(path, s.Router)
	if err != nil {
		// todo: [fix] [A002] Finish the Blunder package and throw an error
	}
}
