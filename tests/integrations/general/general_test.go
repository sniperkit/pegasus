package general_test

import (
	"github.com/cpapidas/pegasus/netamqp"
	"github.com/cpapidas/pegasus/netgrpc"
	"github.com/cpapidas/pegasus/nethttp"
	"github.com/cpapidas/pegasus/peg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneral_allProtocols(t *testing.T) {

	var skipAMQP = false

	defer func() {
		if r := recover(); r != nil {
			skipAMQP = true
		}
	}()

	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0
	serverAMQP := netamqp.NewServer()
	serverAMQP.Serve("amqp://guest:guest@localhost:5672/")

	if skipAMQP {
		t.Skip("You need a RabbitMQ server to continue")
		return
	}

	var handler = func(channel *peg.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := peg.NewOptions().Unmarshal(receive.Options)

		replyOptions := peg.NewOptions()

		// RabbitMQ does not send back any response so we have to do the assertions inside handler
		if options.GetHeader("Custom") == "" || receive.Body == nil {
			panic("Header Custom and Body have to be set")
		}

		replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
		replyOptions.SetHeader("name", options.GetParam("name")+" response")
		replyOptions.SetHeader("id", options.GetParam("id")+" response")

		responseBody := string(receive.Body) + " response"

		// Create the new payload
		payload := peg.BuildPayload([]byte(responseBody), replyOptions.Marshal())

		// Send it back
		channel.Send(payload)

	}

	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	serverHTTP := nethttp.NewServer(nil)
	serverGRPC := netgrpc.NewServer(nil)
	serverAMQP = netamqp.NewServer()

	serverAMQP.Serve("amqp://guest:guest@localhost:5672/")

	serverHTTP.Listen(nethttp.SetConf("/hello", nethttp.Put), handler, nil)
	serverGRPC.Listen(netgrpc.SetConf("/hello"), handler, nil)
	serverAMQP.Listen(netamqp.SetConf("/hello"), handler, nil)

	serverHTTP.Serve("localhost:7001")
	serverGRPC.Serve("localhost:50051")

	// Create a payload
	options := peg.NewOptions()

	options.SetHeader("Custom", "header-value")

	payload := peg.BuildPayload([]byte("foo"), options.Marshal())

	// Send the payload
	response, err := nethttp.NewClient("http://localhost:7001/").
		Send(nethttp.SetConf("hello?name=christos", nethttp.Put), payload)

	replyOptions := peg.NewOptions().Unmarshal(response.Options)

	// Should not throw an error
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err, "Should be nil")

	// The response should have the following values
	assert.Equal(t, "foo response", string(response.Body),
		"Should be equals to foo response")
	assert.Equal(t, "header-value response", replyOptions.GetHeader("Custom"),
		"Should be equals to header-value response")

	// Should returns the param name
	assert.Equal(t, "christos response", replyOptions.GetHeader("Name"),
		"Should be equals to christos response")

	client, _ := netamqp.NewClient("amqp://guest:guest@localhost:5672/")

	// Should not throw panic
	options = peg.NewOptions()
	options.SetHeader("Custom", "bar")
	payload = peg.BuildPayload([]byte("foo"), options.Marshal())

	assert.NotPanics(t, func() { client.Send(netamqp.SetConf("/hello"), payload) },
		"Should not panics")

	// Create a payload
	options = peg.NewOptions()

	options.SetHeader("Custom", "header-value")

	payload = peg.BuildPayload([]byte("foo"), options.Marshal())

	// Send the payload
	response, err = netgrpc.NewClient("localhost:50051").
		Send(netgrpc.SetConf("/hello"), payload)

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
}
