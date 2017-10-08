package nethttp

import (
	"github.com/cpapidas/pegasus/network"
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/cpapidas/pegasus/helpers"
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

	path := conf[0]
	method := conf[1]

	var body []byte
	if method != Get.String() {
		body = payload.Body
	}

	// Create a request
	request, err := NewRequest(method, path, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	httpOptions := network.BuildOptions(payload.Options)
	headers := httpOptions.GetHeaders()

	request.Header.Set("Content-Type", "application/json")

	if headers != nil {
		//Set the HTTP request headers
		for key, value := range headers {
			if helpers.IsHTTPValidHeader(key) {
				request.Header.Set(key, value)
			}
		}
	}

	// Send the request
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	// Build the options for the response
	responseOptions := network.NewOptions()

	// Set the HTTP request headers
	for key, value := range response.Header {
		if helpers.IsHTTPValidHeader(key) {
			responseOptions.SetHeader(key, strings.Join(value, ","))
		}
	}

	// Close the body
	defer response.Body.Close()

	// Get get body content
	content, err := ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return network.NewPayload(content, responseOptions.Marshal()), nil
}

// Close terminal the current connection.
func (Client) Close() error{
	return nil
}
