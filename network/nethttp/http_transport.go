package tranhttp

import (
	"bytes"
	"errors"
	"bitbucket.org/code_horse/pegasus/network"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Doc in transport interface
type Transporter struct {
	// The body is the payload which we sent or we receive
	Body []byte

	// The mux router
	Router *mux.Router
}

func (t *Transporter) Send(properties *network.Properties, options *network.Options, payload []byte) (*network.Payload, error) {

	// Create a new http client
	client := &http.Client{}

	httpProperties := NewProperties().BuildProperties(properties)
	httpMethod := httpProperties.GetMethod()

	if httpMethod == "" {
		return nil, errors.New("network.properties[\"METHOD\"][\"VALUE\"] is empty")
	}

	// Create a request
	request, err := http.NewRequest(httpMethod, httpProperties.GetPath(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	httpOptions := NewOptions().Build(options)
	headers := httpOptions.GetHeaders()

	if headers != nil {
		//Set the HTTP request headers
		for key, value := range headers {
			request.Header.Set(key, value)
		}
	}

	// Send the request
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// Build the options for the response
	responseOptions := NewOptions()

	//Set the HTTP request headers
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

	return &network.Payload{Body: content, Options: responseOptions.Marshal()}, nil
}

func (t *Transporter) Listen(properties *network.Properties, handler network.Handler, middleware network.Middleware) {

	httpProperties := NewProperties().BuildProperties(properties)

	t.Router.HandleFunc(httpProperties.GetPath(), func(w http.ResponseWriter, r *http.Request) {

		requestOptions := NewOptions()

		requestOptions.SetHeaders(t.setHeaders(r.Header))

		requestOptions.SetQueryParams(t.setQueryParams(r.URL.Query()))

		// Parse and set the body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		channel := network.NewChannel(2000)

		requestPayload := network.BuildPayload(body, requestOptions.Marshal())

		channel.Send(requestPayload)

		if middleware != nil {
			middleware(handler, channel)
		} else {
			handler(channel)
		}

		responsePayload := channel.Receive()

		responseOptions := NewOptions().BuildFromBytes(responsePayload.Options)
		responseHeaders := responseOptions.GetHeaders()

		// Set the http headers
		if responseHeaders != nil {
			for key, value := range responseHeaders {
				w.Header().Set(key, value)
			}
		}

		// Return the response
		w.Write(responsePayload.Body)
		//w.Write([]byte("hello"))

	}).Methods(httpProperties.GetMethod())
}

// Set a map of strings keys and strings values of given http headers. Receives
// a http header object and returns a map object map[string]string.
func (t *Transporter) setHeaders(headers http.Header) map[string]string {
	mapper := make(map[string]string)
	for key, value := range headers {
		mapper[key] = strings.Join(value, ",")
	}
	return mapper
}

// Set a map of strings keys and strings values of given url query params. Receives
// a http header object and returns a map object map[string]string.
func (t *Transporter) setQueryParams(params url.Values) map[string]string {
	mapper := make(map[string]string)
	for key, value := range params {
		mapper[key] = strings.Join(value, ",")
	}
	return mapper
}
