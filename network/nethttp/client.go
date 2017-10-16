package nethttp

import (
	"bytes"
	"github.com/cpapidas/pegasus/helpers"
	"github.com/cpapidas/pegasus/network"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client interface describes the protocols client model. Client keeps the connections open for each protocol.
// In order to close the connection the Close function should be used. Send function is used to send data
// to other servers.
type Client struct {
	httpClient IHTTPClient
}

// IHTTPClient interface describe the http.Client struct
type IHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewRequest http.NewRequest
var NewRequest = http.NewRequest

// ReadAll ioutil.ReadAll
var ReadAll = ioutil.ReadAll

// NewClient generates and returns a Client object.
var NewClient = func(httpClient IHTTPClient) network.Client {

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{httpClient: httpClient}
}

// Send function sends a payload to other servers. It gets the string path which is the unique id and the payload
// object. The path may use other function ir order to generate the format for each provider.
func (c Client) Send(conf []string, payload network.Payload) (*network.Payload, error) {

	httpOptions := network.BuildOptions(payload.Options)

	// Create a request
	request, err := c.createRequest(conf[0], conf[1], payload.Body)
	if err != nil {
		return nil, err
	}

	c.setRequestParams(request, httpOptions)
	c.setRequestHeaders(request, httpOptions)

	// Send the request
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	responseOptions := network.NewOptions()
	c.setResponseHeaders(response, responseOptions)

	defer response.Body.Close()

	content, err := ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return network.NewPayload(content, responseOptions.Marshal()), nil
}

// Close terminal the current connection.
func (Client) Close() error {
	return nil
}

// createRequest creates a new http request object
func (Client) createRequest(path string, method string, body []byte) (*http.Request, error) {
	if method == Get.String() {
		body = nil
	}
	return NewRequest(method, path, bytes.NewReader(body))
}

// setRequestParams sets the request parameter
func (Client) setRequestParams(request *http.Request, options *network.Options) {
	q := request.URL.Query()
	for k, v := range options.GetParams() {
		q.Add(k, v)
	}
	request.URL.RawQuery = q.Encode()
}

// setRequestHeaders sets the request headers
func (Client) setRequestHeaders(request *http.Request, options *network.Options) {
	request.Header.Set("Content-Type", "application/json")
	headers := options.GetHeaders()
	if headers != nil {
		for key, value := range headers {
			if helpers.IsHTTPValidHeader(key) {
				request.Header.Set(key, value)
			}
		}
	}
}

// setResponseHeaders set the response headers
func (Client) setResponseHeaders(response *http.Response, options *network.Options) {
	for key, value := range response.Header {
		if helpers.IsHTTPValidHeader(key) {
			options.SetHeader(key, strings.Join(value, ","))
		}
	}
}
