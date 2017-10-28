package nethttp

import (
	"github.com/cpapidas/pegasus/peg"

	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ListenAndServe start the server
var ListenAndServe = http.ListenAndServe

// HandleFunction handle the http request
var HandleFunction = http.HandleFunc

// Server implements the peg.Server
// Server struct is responsible for http server. It manages connections and configuration might be needed in order to
// ensure that the http server works properly
type Server struct {
}

// NewServer is a constructor of Server struct. It initializes and returns a Server object.
var NewServer = func() peg.Server {
	return &Server{}
}

// SetConf gets a path as parameter and returns an array. It is used for Server.Listen.
func SetConf(path string, method Method) []string {
	return []string{path, method.String()}
}

// Serve function starts the server for a specific part and port.
func (s Server) Serve(path string) {
	go func() { ListenAndServe(path, nil) }()
}

// Listen function creates a handler for a specific endpoint. It gets the path string unique key, the handler
// which is a function and the middleware which also is a function.
func (s Server) Listen(conf []string, handler peg.Handler, middleware peg.Middleware) {
	HandleFunction(conf[0], func(w http.ResponseWriter, r *http.Request) {
		options := s.setRequestOptions(r)

		body, err := ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		channel := peg.NewChannel(1)

		requestPayload := peg.BuildPayload(body, options.Marshal())
		channel.Send(requestPayload)
		s.callHandler(handler, middleware, channel)

		s.response(channel, w)
	})
}

// setRequestOptions sets and return a peg.Option object.
func (s Server) setRequestOptions(r *http.Request) *peg.Options {
	options := peg.NewOptions()
	options.SetHeaders(s.setRequestHeaders(r.Header))

	options.SetParams(s.setQueryParams(r.URL.Query()))
	return options
}

// response sends the response via http
func (s Server) response(channel *peg.Channel, w http.ResponseWriter) {
	responsePayload := channel.Receive()
	responseOptions := peg.BuildOptions(responsePayload.Options)
	s.setResponseHeaders(responseOptions, w)
	w.Write(responsePayload.Body)
}

// setRequestHeaders sets a map of string keys and string values of given http headers. Receives
// an http header object and returns a map object map[string]string.
func (Server) setRequestHeaders(headers http.Header) map[string]string {
	mapper := make(map[string]string)
	for key, value := range headers {
		if peg.IsHTTPValidHeader(key) {
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

// callHandler call the http handler. If the middleware para is not nil will call only the middleware.
func (Server) callHandler(handler peg.Handler, middleware peg.Middleware, channel *peg.Channel) {
	if middleware != nil {
		middleware(handler, channel)
	} else {
		handler(channel)
	}
}

// setResponseHeaders sets the response headers
func (s Server) setResponseHeaders(options *peg.Options, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	responseHeaders := options.GetHeaders()
	if responseHeaders != nil {
		for key, value := range responseHeaders {
			if peg.IsHTTPValidHeader(key) {
				w.Header().Set(key, value)
			}
		}
	}
	s.setResponseStatus(w, options.GetHeader("Status"))
}
