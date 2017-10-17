package nethttp

import (
	"github.com/cpapidas/pegasus/helpers"
	"github.com/cpapidas/pegasus/network"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Router interface for mux.Router
type Router interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
}

// Server implements the network.Server
// Server struct is responsible for http server. It manages connections and configuration might be needed in order to
// ensure that the http server works properly
type Server struct {
	// Router is responsible for handler and middleware
	Router Router
}

// NewServer is a constructor of Server struct. It initializes and returns a Server object. It gets a *mux.Router as
// parameter, if the router parameter is nil it will generate a new router and assign it to the object.
var NewServer = func(router Router) network.Server {

	if router == nil {
		router = mux.NewRouter()
	}

	return &Server{Router: router}
}

// SetConf gets a path as parameter and returns an array. It is used for Server.Listen.
func SetConf(path string, method Method) []string {
	return []string{path, method.String()}
}

// Serve function starts the server for a specific part and port.
func (s Server) Serve(path string) {
	go func() {
		err := http.ListenAndServe(path, s.Router.(*mux.Router))
		if err != nil {
			panic(err)
		}
	}()
}

// Listen function creates a handler for a specific endpoint. It gets the path string unique key, the handler
// which is a function and the middleware which also is a function.
func (s Server) Listen(paths []string, handler network.Handler, middleware network.Middleware) {
	s.Router.HandleFunc(paths[0], func(w http.ResponseWriter, r *http.Request) {
		options := s.setRequestOptions(r)

		s.setPathVariables(r, options)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		channel := network.NewChannel(1)

		requestPayload := network.BuildPayload(body, options.Marshal())
		channel.Send(requestPayload)
		s.callHandler(handler, middleware, channel)

		s.response(channel, w)

	}).Methods(paths[1])
}

// setRequestOptions sets and return a network.Option object.
func (s Server) setRequestOptions(r *http.Request) *network.Options {
	options := network.NewOptions()
	options.SetHeaders(s.setRequestHeaders(r.Header))
	options.SetParams(s.setQueryParams(r.URL.Query()))

	return options
}

// response sends the response via http
func (s Server) response(channel *network.Channel, w http.ResponseWriter) {
	responsePayload := channel.Receive()
	responseOptions := network.BuildOptions(responsePayload.Options)
	s.setResponseHeaders(responseOptions, w)
	w.Write(responsePayload.Body)
}

// setRequestHeaders sets a map of string keys and string values of given http headers. Receives
// an http header object and returns a map object map[string]string.
func (Server) setRequestHeaders(headers http.Header) map[string]string {
	mapper := make(map[string]string)
	for key, value := range headers {
		if helpers.IsHTTPValidHeader(key) {
			mapper[key] = strings.Join(value, ",")
		}
	}
	return mapper
}

// setQueryParams sets a map of string keys and string values of given url query params. Receives
// an http header object and returns a map object map[string]string.
func (Server) setQueryParams(params url.Values) map[string]string {
	mapper := make(map[string]string)
	for key, value := range params {
		mapper[key] = strings.Join(value, ",")
	}
	return mapper
}

// setResponseStatus sets the response status. If the header with Status is set then pass the header value.
func (Server) setResponseStatus(w http.ResponseWriter, status string) {
	if status := status; status != "" {
		s, _ := strconv.Atoi(status)
		w.WriteHeader(s)
	}
}

// setPathVariables sets the path variables
func (Server) setPathVariables(r *http.Request, options *network.Options) {
	for pathKey, pathVar := range mux.Vars(r) {
		options.SetParam(pathKey, pathVar)
	}
}

// callHandler call the http handler. If the middleware para is not nil will call only the middleware.
func (Server) callHandler(handler network.Handler, middleware network.Middleware, channel *network.Channel) {
	if middleware != nil {
		middleware(handler, channel)
	} else {
		handler(channel)
	}
}

// setResponseHeaders sets the response headers
func (s Server) setResponseHeaders(options *network.Options, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	responseHeaders := options.GetHeaders()
	if responseHeaders != nil {
		for key, value := range responseHeaders {
			if helpers.IsHTTPValidHeader(key) {
				w.Header().Set(key, value)
			}
		}
	}
	s.setResponseStatus(w, options.GetHeader("Status"))
}
