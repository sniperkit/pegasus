package nethttp_test

import (
	"bytes"
	"errors"
	"github.com/cpapidas/pegasus/nethttp"
	"github.com/cpapidas/pegasus/peg"
	"github.com/cpapidas/pegasus/tests/mocks/mhttp"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {

	client := nethttp.NewClient("http://localhost")

	// Should not be nil
	assert.NotNil(t, client, "Should not return a nil client object")

	// Should be type of *nethttp.Client
	assert.Equal(t, "<*nethttp.Client Value>", reflect.ValueOf(client).String(),
		"Should be type of <*nethttp.Client Value>")
}

func TestClient_Send(t *testing.T) {

	called := false

	// Create a mock object
	mockHTTPClient := &mhttp.MockHTTPClient{}

	// Mock the methods
	mockHTTPClient.DoMock = func(req *http.Request) (*http.Response, error) {

		called = true

		// Should have the right parameters
		assert.Equal(t, []string{"sample"}, req.Header["Custom-Sample"],
			"Should be equals to [sample]")
		assert.Equal(t, []string{"sample"}, req.Header["Hp-Sample"],
			"Should be equals to [sample]")
		assert.Nil(t, req.Header["Mq-Sample"], "Should be nil")
		assert.Nil(t, req.Header["Gr-Sample"], "Should be nil")

		b, err := ioutil.ReadAll(req.Body)

		assert.Nil(t, err, "Should not be nil")
		assert.Equal(t, "content", string(b), "Should be equals to content")

		response := &http.Response{}
		response.Header = make(map[string][]string)
		response.Header["Custom-Sample-Reply"] = []string{"sample reply"}
		response.Header["HP-Sample-Reply"] = []string{"sample reply"}
		response.Header["MQ-Sample-Reply"] = []string{"sample reply"}
		response.Header["GR-Sample-Reply"] = []string{"sample reply"}

		response.Body = ioutil.NopCloser(bytes.NewReader([]byte("content reply")))
		response.StatusCode = 200

		return response, nil
	}

	requestOptions := peg.NewOptions()
	requestOptions.SetHeader("Custom-Sample", "sample")
	requestOptions.SetHeader("HP-Sample", "sample")
	requestOptions.SetHeader("MQ-Sample", "sample")
	requestOptions.SetHeader("GR-Sample", "sample")

	// Pass the mock object and create a client
	client := nethttp.NewClient("", mockHTTPClient)

	payload := peg.BuildPayload([]byte("content"), requestOptions.Marshal())

	reply, err := client.Send(nethttp.SetConf("/whatever", "POST"), payload)

	replyOptions := peg.NewOptions().Unmarshal(reply.Options)

	// Should return a nil error
	assert.Nil(t, err, "Should not returns error")

	// Should call the Do function
	assert.True(t, called, "Should be call the handler")

	// Should return valid parameters
	assert.Equal(t, "sample reply", replyOptions.GetHeader("Custom-Sample-Reply"),
		"Should be equals to sample reply")
	assert.Equal(t, "sample reply", replyOptions.GetHeader("HP-Sample-Reply"),
		"Should be equals to sample reply")
	assert.Empty(t, replyOptions.GetHeader("MQ-Sample-Reply"), "Should be empty")
	assert.Empty(t, replyOptions.GetHeader("GR-Sample-Reply"), "Should be empty")
	assert.Equal(t, "200", replyOptions.GetHeader("Status"))
	assert.Equal(t, string(reply.Body), "content reply", "Should be equals to content reply")
}

func TestClient_Send_getRequest(t *testing.T) {

	called := false

	mockHTTPClient := &mhttp.MockHTTPClient{}

	mockHTTPClient.DoMock = func(req *http.Request) (*http.Response, error) {

		called = true

		// Should have the right parameters
		assert.Equal(t, []string{"sample"}, req.Header["Custom-Sample"],
			"Should be equals to [sample]")
		assert.Equal(t, []string{"sample"}, req.Header["Hp-Sample"],
			"Should be equals to [sample]")
		assert.Nil(t, req.Header["Mq-Sample"], "Should be nil")
		assert.Nil(t, req.Header["Gr-Sample"], "Should be nil")

		b, err := ioutil.ReadAll(req.Body)

		assert.Nil(t, err, "Should not be nil")
		assert.Empty(t, string(b), "Should be empty")

		assert.Equal(t, "bar", req.URL.Query().Get("foo"),
			"Should be equals to bar")

		response := &http.Response{}
		response.Header = make(map[string][]string)
		response.Header["Custom-Sample-Reply"] = []string{"sample reply"}
		response.Header["HP-Sample-Reply"] = []string{"sample reply"}
		response.Header["MQ-Sample-Reply"] = []string{"sample reply"}
		response.Header["GR-Sample-Reply"] = []string{"sample reply"}

		response.Body = ioutil.NopCloser(bytes.NewReader([]byte("content reply")))

		return response, nil
	}

	requestOptions := peg.NewOptions()
	requestOptions.SetHeader("Custom-Sample", "sample")
	requestOptions.SetHeader("HP-Sample", "sample")
	requestOptions.SetHeader("MQ-Sample", "sample")
	requestOptions.SetHeader("GR-Sample", "sample")

	requestOptions.SetParam("foo", "bar")

	client := nethttp.NewClient("", mockHTTPClient)

	payload := peg.BuildPayload([]byte("content"), requestOptions.Marshal())

	reply, err := client.Send(nethttp.SetConf("/whatever", nethttp.Get), payload)

	replyOptions := peg.NewOptions().Unmarshal(reply.Options)

	// Should return a nil error
	assert.Nil(t, err, "Should not returns error")

	// Should call the Do function
	assert.True(t, called, "Should be call the handler")

	// Should return valid parameters
	assert.Equal(t, "sample reply", replyOptions.GetHeader("Custom-Sample-Reply"),
		"Should be equals to sample reply")
	assert.Equal(t, "sample reply", replyOptions.GetHeader("HP-Sample-Reply"),
		"Should be equals to sample reply")
	assert.Empty(t, replyOptions.GetHeader("MQ-Sample-Reply"), "Should be empty")
	assert.Empty(t, replyOptions.GetHeader("GR-Sample-Reply"), "Should be empty")
	assert.Equal(t, string(reply.Body), "content reply", "Should be equals to content reply")
}

func TestClient_Send_requestFailure(t *testing.T) {

	nethttp.NewRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("error")
	}

	client := nethttp.NewClient("http://localhost")

	payload, err := client.Send(nethttp.SetConf("what", "ever"), peg.Payload{})

	// Should return nil payload
	assert.Nil(t, payload, "Should be nil")

	// Should return an error
	assert.NotNil(t, err, "Should not be nil")

	nethttp.NewRequest = http.NewRequest
}

func TestClient_Send_doFailure(t *testing.T) {

	called := false

	mockHTTPClient := &mhttp.MockHTTPClient{}

	mockHTTPClient.DoMock = func(req *http.Request) (*http.Response, error) {
		called = true
		return nil, errors.New("error")
	}

	client := nethttp.NewClient("", mockHTTPClient)

	reply, err := client.Send(nethttp.SetConf("/whatever", "POST"), peg.Payload{})

	// Should call the Do function
	assert.True(t, called, "Should call the handler")

	// Should return a nil reply
	assert.Nil(t, reply, "Should be nil")

	// Should an error
	assert.NotNil(t, err, "Should not be nil")
}

func TestClient_Send_readAllFailure(t *testing.T) {

	called := false

	nethttp.ReadAll = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("error")
	}

	mockHTTPClient := &mhttp.MockHTTPClient{}

	mockHTTPClient.DoMock = func(req *http.Request) (*http.Response, error) {
		called = true
		response := &http.Response{}
		response.Body = ioutil.NopCloser(bytes.NewReader([]byte("content reply")))
		return response, nil
	}

	client := nethttp.NewClient("", mockHTTPClient)

	reply, err := client.Send(nethttp.SetConf("/whatever", "POST"), peg.Payload{})

	// Should call the Do function
	assert.True(t, called, "Should be true")

	// Should return a nil reply
	assert.Nil(t, reply, "Should be nil")

	// Should an error
	assert.NotNil(t, err, "Should not be nil")

	nethttp.ReadAll = ioutil.ReadAll
}

func TestClient_Close(t *testing.T) {
	client := nethttp.NewClient("http://localhost")
	// Should an error, because this function is not implemented for nethttp provider
	assert.NotNil(t, client.Close(), "Should not be nil")
}
