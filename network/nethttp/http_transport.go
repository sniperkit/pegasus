package nethttp
//
//import (
//	"bytes"
//	"errors"
//	"bitbucket.org/code_horse/pegasus/network"
//	"github.com/gorilla/mux"
//	"io/ioutil"
//	"net/http"
//	"net/url"
//	"strings"
//)
//
//// Doc in transport interface
//type Transporter struct {
//	// The body is the payload which we sent or we receive
//	Body []byte
//
//	// The mux router
//	Router *mux.Router
//}
//
//func (t *Transporter) Send(properties *network.Properties, options *network.Options, payload []byte) (*network.Payload, error) {
//
//	// Create a new http client
//	client := &http.Client{}
//
//	httpProperties := NewProperties().BuildProperties(properties)
//	httpMethod := httpProperties.GetMethod()
//
//	if httpMethod == "" {
//		return nil, errors.New("network.properties[\"METHOD\"][\"VALUE\"] is empty")
//	}
//
//	// Create a request
//	request, err := http.NewRequest(httpMethod, httpProperties.GetPath(), bytes.NewReader(payload))
//	if err != nil {
//		return nil, err
//	}
//
//	httpOptions := NewOptions().Build(options)
//	headers := httpOptions.GetHeaders()
//
//	if headers != nil {
//		//Set the HTTP request headers
//		for key, value := range headers {
//			request.Header.Set(key, value)
//		}
//	}
//
//	// Send the request
//	response, err := client.Do(request)
//	if err != nil {
//		return nil, err
//	}
//
//	// Build the options for the response
//	responseOptions := NewOptions()
//
//	//Set the HTTP request headers
//	for key, value := range response.Header {
//		responseOptions.SetHeader(key, strings.Join(value, ","))
//	}
//
//	// Close the body
//	defer response.Body.Close()
//
//	// Get get body content
//	content, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	return &network.Payload{Body: content, Options: responseOptions.Marshal()}, nil
//}