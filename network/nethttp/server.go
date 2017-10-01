package nethttp

import (
	"bitbucket.org/code_horse/pegasus/network"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Server implements the network.Server
// Server struct is responsible for http server. It manages connections and configuration may needed in order to ensure
// that the http server works properly
type Server struct {
	// Router is responsible for handler and middleware
	Router *mux.Router
}

// NewServer is constructor of Server struct. It initialize and return a Server object. It get a *mux.Router as
// parameter, if the router parameter is nil it will generate a new router and assign it to the object.
var NewServer = func(router *mux.Router) network.Server {

	if router == nil {
		router = mux.NewRouter()
	}

	return &Server{Router: router}
}

// SetPath get a path as parameter and return an array. It use for Server.Listen.
func SetPath(path string, method Method) []string {
	return []string{path, method.String()}
}

// Serve function start the server for a specific part and port
func (s Server) Serve(path string) {
	go func() {
		err := http.ListenAndServe(path, s.Router)
		if err != nil {
			panic(err)
		}
	}()
}

// Listen function creates a handler for a specific endpoint. It gets the path string unique key, the handler
// which is a function and the middleware which also is a function.
func (s Server) Listen(paths []string, handler network.Handler, middleware network.Middleware) {

	path := paths[0]
	method := paths[1]

	s.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		options := network.NewOptions()

		// Set the headers
		options.SetHeaders(s.setHeaders(r.Header))

		// Set the query params
		options.SetParams(s.setQueryParams(r.URL.Query()))

		// Set path variables
		for pathKey, pathVar := range mux.Vars(r) {
			options.SetParam(pathKey, pathVar)
		}

		// Parse and set the body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// todo: [fix] [A006] Return the error at handler
			panic(err.Error())
		}

		//todo: [fix] [A007] Make channel infinity
		channel := network.NewChannel(1)

		// Get the payload
		requestPayload := network.BuildPayload(body, options.Marshal())

		// Send the payload
		channel.Send(requestPayload)

		// Start the handler
		if middleware != nil {
			middleware(handler, channel)
		} else {
			handler(channel)
		}

		// Receive the data from channel
		responsePayload := channel.Receive()

		// Build the options form bytes to nethttp.Option object
		responseOptions := network.BuildOptions(responsePayload.Options)

		// Get the headers from options
		responseHeaders := responseOptions.GetHeaders()

		// Set the http headers
		if responseHeaders != nil {
			for key, value := range responseHeaders {
				w.Header().Set(key, value)
			}
		}

		// If the header with Status is set then pass the header value
		if status := responseOptions.GetHeader("Status"); status != "" {
			s, _ := strconv.Atoi(status)
			w.WriteHeader(s)
		}

		w.Write(responsePayload.Body)

	}).Methods(method)

}

// setHeaders sets a map of strings keys and strings values of given http headers. Receives
// a http header object and returns a map object map[string]string.
func (Server) setHeaders(headers http.Header) map[string]string {
	mapper := make(map[string]string)
	for key, value := range headers {
		mapper[key] = strings.Join(value, ",")
	}
	return mapper
}

// setQueryParams sets a map of strings keys and strings values of given url query params. Receives
// a http header object and returns a map object map[string]string.
func (Server) setQueryParams(params url.Values) map[string]string {
	mapper := make(map[string]string)
	for key, value := range params {
		mapper[key] = strings.Join(value, ",")
	}
	return mapper
}
