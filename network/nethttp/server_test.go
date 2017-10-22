package nethttp_test

import (
	"github.com/cpapidas/pegasus/network/nethttp"

	"bytes"
	"github.com/cpapidas/pegasus/network"
	"github.com/cpapidas/pegasus/tests/mocks/mhttp"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"reflect"
	"io"
	"errors"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestNewServer(t *testing.T) {
	// Should not be nil
	server := nethttp.NewServer(nil)
	assert.NotNil(t, server, "Should not be nil")

	// Should be type of *Server
	server = nethttp.NewServer(nil)
	assert.Equal(t, "<*nethttp.Server Value>", reflect.ValueOf(server).String(),
		"Should be type of <*nethttp.Server Value>")
}

func TestSetConf(t *testing.T) {
	// Should return an array of given strings
	assert.Equal(t, []string{"foo", "bar"}, nethttp.SetConf("foo", "bar"),
		`Should be equals to []string{"foo", "bar"}`)
}

func TestServer_Serve(t *testing.T) {
	var server network.Server
	nethttp.ListenAndServe = func(addr string, handler http.Handler) error {
		return nil
	}
	server = nethttp.NewServer(nil)
	assert.NotPanics(t, func() { server.Serve("Foo") }, "Should not panics")
}

func TestServer_Listen(t *testing.T) {
	callHandler := false

	// Create a mock router
	router := &mhttp.MockRouter{}

	// Initialize the handler function
	router.HandleFuncMock = func(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {

		w := &mhttp.MockResponseWriter{
			Headers: make(map[string][]string),
		}

		r, _ := http.NewRequest("POST", "anything", bytes.NewReader([]byte("content")))

		r.Header = make(map[string][]string)
		r.Header["Custom-Sample"] = []string{"Sample"}
		r.Header["HP-Sample"] = []string{"Sample"}
		r.Header["GR-Sample"] = []string{"Sample"}
		r.Header["MQ-Sample"] = []string{"Sample"}

		r.Body = ioutil.NopCloser(bytes.NewReader([]byte("content")))

		f(w, r)

		// Should contain the right headers
		assert.Equal(t, []string{"sample"}, w.Headers["Custom-Sample"],
			`Should be equals to []string{"sample"}`)
		assert.Equal(t, []string{"sample"}, w.Headers["Hp-Sample"],
			`Should be equals to []string{"sample"}`)
		assert.Empty(t, w.Headers["Gr-Sample"],
			"Should be empty")
		assert.Empty(t, w.Headers["Mq-Sample"],
			"Should be empty")
		assert.Equal(t, "content reply", string(w.Body),
			"Should be content reply")

		return &mux.Route{}
	}

	server := nethttp.NewServer(router)

	// Initialize the request handler
	var handler = func(channel *network.Channel) {

		callHandler = true

		payload := channel.Receive()

		options := network.NewOptions().Unmarshal(payload.Options)

		// Should return the valid headers
		assert.Equal(t, "Sample", options.GetHeader("Custom-Sample"),
			"Should be equals to Sample")
		assert.Equal(t, "Sample", options.GetHeader("HP-Sample"),
			"Should be equals to Sample")
		assert.Empty(t, options.GetHeader("GR-Sample"), "Should be empty")
		assert.Empty(t, options.GetHeader("MQ-Sample"), "Should be empty")

		replyOptions := network.NewOptions()

		replyOptions.SetHeader("Custom-Sample", "sample")
		replyOptions.SetHeader("HP-Sample", "sample")
		replyOptions.SetHeader("GR-Sample", "sample")
		replyOptions.SetHeader("MQ-Sample", "sample")

		channel.Send(network.BuildPayload([]byte("content reply"), replyOptions.Marshal()))

	}

	server.Listen(nethttp.SetConf("foo", "POST"), handler, nil)

	// Should call the handler
	assert.True(t, callHandler, "Should call the handler")
}

func TestServer_Listen_middleware(t *testing.T) {
	callHandler := false
	callMiddleware := false

	router := &mhttp.MockRouter{}

	// Initialize the handler function
	router.HandleFuncMock = func(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {

		w := &mhttp.MockResponseWriter{
			Headers: make(map[string][]string),
		}

		r, _ := http.NewRequest("PUT", "/url/path", bytes.NewReader([]byte("content")))

		r.Header = make(map[string][]string)
		r.Header["Custom-Sample"] = []string{"Sample"}
		r.Header["HP-Sample"] = []string{"Sample"}
		r.Header["GR-Sample"] = []string{"Sample"}
		r.Header["MQ-Sample"] = []string{"Sample"}

		r.Body = ioutil.NopCloser(bytes.NewReader([]byte("content")))

		q := r.URL.Query()
		q.Add("foo", "bar")
		r.URL.RawQuery = q.Encode()

		f(w, r)

		// Should contain the right headers
		assert.Equal(t, []string{"sample"}, w.Headers["Custom-Sample"],
			`Should be equals to []string{"sample"}`)
		assert.Equal(t, []string{"sample"}, w.Headers["Hp-Sample"],
			`Should be equals to []string{"sample"}`)
		assert.Empty(t, w.Headers["Gr-Sample"],
			"Should be empty")
		assert.Empty(t, w.Headers["Mq-Sample"],
			"Should be empty")
		assert.Equal(t, "content reply", string(w.Body),
			"Should be content reply")
		assert.Equal(t, 201, w.Status, "Should be equals to 201")

		return &mux.Route{}
	}

	// Create a server object with mocked date
	server := nethttp.NewServer(router)

	// Initialize the handler
	var handler = func(channel *network.Channel) {

		callHandler = true

		payload := channel.Receive()

		options := network.NewOptions().Unmarshal(payload.Options)

		// Should return the valid headers
		assert.Equal(t, "sample middleware", options.GetHeader("Custom-Sample"),
			"Should be equals to sample middleware")
		assert.Equal(t, "sample middleware", options.GetHeader("HP-Sample"),
			"Should be equals to sample middleware")

		replyOptions := network.NewOptions()

		replyOptions.SetHeader("Custom-Sample", "sample")
		replyOptions.SetHeader("HP-Sample", "sample")
		replyOptions.SetHeader("GR-Sample", "sample")
		replyOptions.SetHeader("MQ-Sample", "sample")
		replyOptions.SetHeader("Status", "201")

		channel.Send(network.BuildPayload([]byte("content reply"), replyOptions.Marshal()))

	}

	// Middleware router
	var middleware = func(handler network.Handler, channel *network.Channel) {

		callMiddleware = true

		payload := channel.Receive()

		options := network.NewOptions().Unmarshal(payload.Options)

		// Should contains the valid header
		assert.Equal(t, "Sample", options.GetHeader("Custom-Sample"),
			"Should be equals to Sample")
		assert.Equal(t, "Sample", options.GetHeader("HP-Sample"),
			"Should be equals to Sample")
		assert.Empty(t, options.GetHeader("GR-Sample"), "Should be empty")
		assert.Empty(t, options.GetHeader("MQ-Sample"), "Should be empty")

		// Should contains the right parameters
		assert.Equal(t, "bar", options.GetParam("foo"),
			"Should be equals to bar")

		replyOptions := network.NewOptions()

		replyOptions.SetHeader("Custom-Sample", "sample middleware")
		replyOptions.SetHeader("HP-Sample", "sample middleware")

		channel.Send(network.BuildPayload([]byte("content reply middleware"), replyOptions.Marshal()))

		handler(channel)
	}

	server.Listen(nethttp.SetConf("/url/path", "PUT"), handler, middleware)

	// Should call the handler
	assert.True(t, callHandler, "Should call the handler")

	// Should call the middleware
	assert.True(t, callMiddleware, "Should call the middleware")
}

func TestServer_Listen_readAllFailure(t *testing.T) {
	var server network.Server
	router := &mhttp.MockRouter{}
	router.HandleFuncMock = func(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
		w := &mhttp.MockResponseWriter{
			Headers: make(map[string][]string),
		}
		r, _ := http.NewRequest("POST", "anything", bytes.NewReader([]byte("content")))
		f(w, r)

		return &mux.Route{}
	}
	server = nethttp.NewServer(router)
	nethttp.ReadAll = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("error")
	}
	// Should not panic
	assert.Panics(t, func() { server.Listen(nethttp.SetConf("", ""), nil, nil) },
		"Should panics")
}
