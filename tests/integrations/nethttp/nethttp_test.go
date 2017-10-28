package nethttp_test

import (
	"github.com/cpapidas/pegasus/nethttp"
	"github.com/cpapidas/pegasus/peg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNethttp_integration(t *testing.T) {

	var handlerGet = func(channel *peg.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := peg.NewOptions().Unmarshal(receive.Options)

		if options.GetHeader("Content-Type") != "application/json" {
			panic("The header Content-Type should have default value application/json")
		}

		replyOptions := peg.NewOptions()

		replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")

		// Create the new payload
		payload := peg.BuildPayload([]byte(options.GetParam("foo")+" response"), replyOptions.Marshal())

		// Send it back
		channel.Send(payload)
	}

	var handlerPost = func(channel *peg.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := peg.NewOptions().Unmarshal(receive.Options)

		replyOptions := peg.NewOptions()

		replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
		replyOptions.SetHeader("name", options.GetParam("name")+" response")

		responseBody := string(receive.Body) + " response"

		// Create the new payload
		payload := peg.BuildPayload([]byte(responseBody), replyOptions.Marshal())

		// Send it back
		channel.Send(payload)
	}

	var handlerPut = func(channel *peg.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := peg.NewOptions().Unmarshal(receive.Options)

		replyOptions := peg.NewOptions()

		replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
		replyOptions.SetHeader("name", options.GetParam("name")+" response")
		replyOptions.SetHeader("id", options.GetParam("id")+" response")

		responseBody := string(receive.Body) + " response"

		// Create the new payload
		payload := peg.BuildPayload([]byte(responseBody), replyOptions.Marshal())

		// Send it back
		channel.Send(payload)
	}

	var handlerDelete = func(channel *peg.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := peg.NewOptions().Unmarshal(receive.Options)

		replyOptions := peg.NewOptions()

		replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
		replyOptions.SetHeader("name", options.GetParam("name")+" response")
		replyOptions.SetHeader("id", options.GetParam("id")+" response")

		// Create the new payload
		payload := peg.BuildPayload([]byte(string(receive.Body)+" response"), replyOptions.Marshal())

		// Send it back
		channel.Send(payload)
	}

	var middleware = func(handler peg.Handler, channel *peg.Channel) {

		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := peg.NewOptions().Unmarshal(receive.Options)

		options.SetHeader("Custom", options.GetHeader("Custom")+" middleware")

		// Create the new payload
		payload := peg.BuildPayload(nil, options.Marshal())

		// Send it back
		channel.Send(payload)

		handler(channel)
	}

	server := nethttp.NewServer(nil)

	server.Listen(nethttp.SetConf("/http", nethttp.Get), handlerGet, nil)
	server.Listen(nethttp.SetConf("/http", nethttp.Post), handlerPost, nil)
	server.Listen(nethttp.SetConf("/http", nethttp.Put), handlerPut, nil)
	server.Listen(nethttp.SetConf("/http", nethttp.Delete), handlerDelete, nil)

	server.Listen(nethttp.SetConf("/http/middleware", nethttp.Get), handlerGet, middleware)

	server.Serve("localhost:7000")

	// HTTP Server

	// Exchange message via HTTP

	// Should not be nil
	assert.NotNil(t, server, "Should not be nil")

	// Send a GET request
	// Create a payload
	options := peg.NewOptions()

	options.SetHeader("Custom", "header-value")
	options.SetHeader("MQ-Custom", "mq-value")
	options.SetHeader("GR-Custom", "gr-value")

	payload := peg.BuildPayload(nil, options.Marshal())

	// Send the payload
	response, err := nethttp.NewClient("http://localhost:7000/").
		Send(nethttp.SetConf("http?foo=bar", nethttp.Get), payload)

	replyOptions := peg.NewOptions().Unmarshal(response.Options)

	// Should not throw an error
	assert.Nil(t, err, "Should not be nil")

	// Should return a application/json as default Content-Type
	assert.Equal(t, "application/json", replyOptions.GetHeader("Content-Type"),
		"Should be equals to application/json")

	// Should return nil headers for GR-* and MQ-*
	assert.Empty(t, replyOptions.GetHeader("MQ-Custom"), "Should be empty")
	assert.Empty(t, replyOptions.GetHeader("GR-Custom"), "Should be empty")

	// She response should have the following values
	assert.Equal(t, []byte("bar response"), response.Body,
		"Should be equals to bar response")
	assert.Equal(t, "header-value response", replyOptions.GetHeader("Custom"),
		"Should be equals to header-value response")

	// Send a POST request", func()
	// Create a payload
	options = peg.NewOptions()

	options.SetHeader("Custom", "header-value")

	payload = peg.BuildPayload([]byte("foo"), options.Marshal())

	// Send the payload
	response, err = nethttp.NewClient("http://localhost:7000/").
		Send(nethttp.SetConf("http?name=christos", nethttp.Post), payload)

	replyOptions = peg.NewOptions().Unmarshal(response.Options)

	// Should not throw an error
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err, "Should not be nil")

	// The response should have the following values
	assert.Equal(t, []byte("foo response"), response.Body)
	assert.Equal(t, "header-value response", replyOptions.GetHeader("Custom"),
		"Should be equals to header-value response")

	// Should returns the param name
	assert.Equal(t, "christos response", replyOptions.GetHeader("Name"),
		"Should be equals to christos response")

	// Send a PUT request
	// Create a payload
	options = peg.NewOptions()

	options.SetHeader("Custom", "header-value")

	payload = peg.BuildPayload([]byte("foo"), options.Marshal())

	// Send the payload
	response, err = nethttp.NewClient("http://localhost:7000/").
		Send(nethttp.SetConf("http?name=christos", nethttp.Put), payload)

	replyOptions = peg.NewOptions().Unmarshal(response.Options)

	// Should not throw an error
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err, "Should be nil")

	// The response should have the following values
	assert.Equal(t, []byte("foo response"), response.Body, "Should be equals to foo repsonse")
	assert.Equal(t, "header-value response", replyOptions.GetHeader("Custom"),
		"Should be equals to header-value response")

	// Should returns the param name
	assert.Equal(t, "christos response", replyOptions.GetHeader("Name"),
		"Should be equals to christos response")

	// Send a DELETE request
	// Create a payload
	options = peg.NewOptions()

	options.SetHeader("Custom", "header-value")

	payload = peg.BuildPayload([]byte("foo"), options.Marshal())

	// Send the payload
	response, err = nethttp.NewClient("http://localhost:7000/").
		Send(nethttp.SetConf("http?name=christos", nethttp.Delete), payload)

	replyOptions = peg.NewOptions().Unmarshal(response.Options)

	// Should not throw an error
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err, "Should be nil")

	// The response should have the following values
	assert.Equal(t, []byte("foo response"), response.Body, "Should be equals to foo response")
	assert.Equal(t, "header-value response", replyOptions.GetHeader("Custom"),
		"Should be equals to header-value response")

	// Should returns the param name
	assert.Equal(t, "christos response", replyOptions.GetHeader("Name"),
		"Should be equal to christos response")

	// Send a GET middleware request
	// Create a payload
	options = peg.NewOptions()

	options.SetHeader("Custom", "header-value")

	payload = peg.BuildPayload(nil, options.Marshal())

	// Send the payload
	response, err = nethttp.NewClient("http://localhost:7000/").
		Send(nethttp.SetConf("http/middleware?foo=bar", nethttp.Get), payload)

	// Should not throw an error", func() {
	assert.Nil(t, err, "Should be nil")

	// The response should have the following values", func() {
	assert.Equal(t, []byte("bar response"), response.Body, "Should be bar response")
	options = peg.NewOptions().Unmarshal(response.Options)
	assert.Equal(t, "header-value middleware response", options.GetHeader("Custom"),
		"Should be equals to header-value middleware response")
}
