package nethttp_test

import (
	"github.com/cpapidas/pegasus/nethttp"

	"bytes"
	"errors"
	"github.com/cpapidas/pegasus/peg"
	"github.com/cpapidas/pegasus/tests/mocks/mhttp"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"

	"testing"
)

func SetDefaults() {
	nethttp.HandleFunction = http.HandleFunc
}

func TestNewServer(t *testing.T) {
	SetDefaults()
	// Should not be nil
	server := nethttp.NewServer()
	assert.NotNil(t, server, "Should not be nil")

	// Should be type of *Server
	server = nethttp.NewServer()
	assert.Equal(t, "<*nethttp.Server Value>", reflect.ValueOf(server).String(),
		"Should be type of <*nethttp.Server Value>")
}

func TestSetConf(t *testing.T) {
	SetDefaults()
	// Should return an array of given strings
	assert.Equal(t, []string{"foo", "bar"}, nethttp.SetConf("foo", "bar"),
		`Should be equals to []string{"foo", "bar"}`)
}

func TestServer_Serve(t *testing.T) {
	SetDefaults()
	var server peg.Server
	nethttp.ListenAndServe = func(addr string, handler http.Handler) error {
		return nil
	}
	server = nethttp.NewServer()
	assert.NotPanics(t, func() { server.Serve("Foo") }, "Should not panics")
}

func TestServer_Listen(t *testing.T) {
	SetDefaults()
	callHandler := false

	// Initialize the handler function
	nethttp.HandleFunction = func(path string, f func(http.ResponseWriter, *http.Request)) {

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
	}

	server := nethttp.NewServer()

	// Initialize the request handler
	var handler = func(channel *peg.Channel) {

		callHandler = true

		payload := channel.Receive()

		options := peg.NewOptions().Unmarshal(payload.Options)

		// Should return the valid headers
		assert.Equal(t, "Sample", options.GetHeader("Custom-Sample"),
			"Should be equals to Sample")
		assert.Equal(t, "Sample", options.GetHeader("HP-Sample"),
			"Should be equals to Sample")
		assert.Empty(t, options.GetHeader("GR-Sample"), "Should be empty")
		assert.Empty(t, options.GetHeader("MQ-Sample"), "Should be empty")

		replyOptions := peg.NewOptions()

		replyOptions.SetHeader("Custom-Sample", "sample")
		replyOptions.SetHeader("HP-Sample", "sample")
		replyOptions.SetHeader("GR-Sample", "sample")
		replyOptions.SetHeader("MQ-Sample", "sample")

		channel.Send(peg.BuildPayload([]byte("content reply"), replyOptions.Marshal()))

	}

	server.Listen(nethttp.SetConf("foo", "POST"), handler, nil)

	// Should call the handler
	assert.True(t, callHandler, "Should call the handler")
}

func TestServer_Listen_middleware(t *testing.T) {
	SetDefaults()
	callHandler := false
	callMiddleware := false

	// Initialize the handler function
	nethttp.HandleFunction = func(path string, f func(http.ResponseWriter, *http.Request)) {

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
	}

	// Create a server object with mocked date
	server := nethttp.NewServer()

	// Initialize the handler
	var handler = func(channel *peg.Channel) {

		callHandler = true

		payload := channel.Receive()

		options := peg.NewOptions().Unmarshal(payload.Options)

		// Should return the valid headers
		assert.Equal(t, "sample middleware", options.GetHeader("Custom-Sample"),
			"Should be equals to sample middleware")
		assert.Equal(t, "sample middleware", options.GetHeader("HP-Sample"),
			"Should be equals to sample middleware")

		replyOptions := peg.NewOptions()

		replyOptions.SetHeader("Custom-Sample", "sample")
		replyOptions.SetHeader("HP-Sample", "sample")
		replyOptions.SetHeader("GR-Sample", "sample")
		replyOptions.SetHeader("MQ-Sample", "sample")
		replyOptions.SetHeader("Status", "201")

		channel.Send(peg.BuildPayload([]byte("content reply"), replyOptions.Marshal()))

	}

	// Middleware router
	var middleware = func(handler peg.Handler, channel *peg.Channel) {

		callMiddleware = true

		payload := channel.Receive()

		options := peg.NewOptions().Unmarshal(payload.Options)

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

		replyOptions := peg.NewOptions()

		replyOptions.SetHeader("Custom-Sample", "sample middleware")
		replyOptions.SetHeader("HP-Sample", "sample middleware")

		channel.Send(peg.BuildPayload([]byte("content reply middleware"), replyOptions.Marshal()))

		handler(channel)
	}

	server.Listen(nethttp.SetConf("/url/path", "PUT"), handler, middleware)

	// Should call the handler
	assert.True(t, callHandler, "Should call the handler")

	// Should call the middleware
	assert.True(t, callMiddleware, "Should call the middleware")
}

func TestServer_Listen_readAllFailure(t *testing.T) {
	SetDefaults()
	var server peg.Server
	nethttp.HandleFunction = func(path string, f func(http.ResponseWriter, *http.Request)) {
		w := &mhttp.MockResponseWriter{
			Headers: make(map[string][]string),
		}
		r, _ := http.NewRequest("POST", "anything", bytes.NewReader([]byte("content")))
		f(w, r)

	}
	server = nethttp.NewServer()
	nethttp.ReadAll = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("error")
	}
	// Should not panic
	assert.Panics(t, func() { server.Listen(nethttp.SetConf("", ""), nil, nil) },
		"Should panics")
}
