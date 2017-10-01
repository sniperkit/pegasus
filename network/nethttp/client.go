package nethttp

import (
	"bitbucket.org/code_horse/pegasus/network"
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client interface describes the protocols client model. Client keep the connections open for each protocol in order
// close the connection should use the Close function. Send function used to send data to other servers.
type Client struct {
	httpClient *http.Client
}

// NewClient generate and return a Client object.
var NewClient = func() network.Client {
	client := &http.Client{}
	return &Client{httpClient: client}
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
	request, err := http.NewRequest(method, path, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	httpOptions := network.BuildOptions(payload.Options)
	headers := httpOptions.GetHeaders()

	//todo: [fix] [A009] Implement one of those by default
	//request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Type", "text/plain")

	if headers != nil {
		//Set the HTTP request headers
		for key, value := range headers {
			request.Header.Set(key, value)
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
		responseOptions.SetHeader(key, strings.Join(value, ","))
	}

	// Close the body
	defer response.Body.Close()

	// Get get body content
	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return network.NewPayload(content, responseOptions.Marshal()), nil
}

// Close terminal the current connection.
func (Client) Close() {
}
