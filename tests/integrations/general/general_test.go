package general_test

import (
	"github.com/cpapidas/pegasus/network"
	"github.com/cpapidas/pegasus/network/netamqp"
	"github.com/cpapidas/pegasus/network/netgrpc"
	"github.com/cpapidas/pegasus/network/nethttp"
	"testing"
	"github.com/stretchr/testify/assert"
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

	var handler = func(channel *network.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := network.NewOptions().Unmarshal(receive.Options)

		replyOptions := network.NewOptions()

		// RabbitMQ does not send back any response so we have to do the assertions inside handler
		if options.GetHeader("Custom") == "" || receive.Body == nil {
			panic("Header Custom and Body have to be set")
		}

		replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
		replyOptions.SetHeader("name", options.GetParam("name")+" response")
		replyOptions.SetHeader("id", options.GetParam("id")+" response")

		responseBody := string(receive.Body) + " response"

		// Create the new payload
		payload := network.BuildPayload([]byte(responseBody), replyOptions.Marshal())

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
	options := network.NewOptions()

	options.SetHeader("Custom", "header-value")

	payload := network.BuildPayload([]byte("foo"), options.Marshal())

	// Send the payload
	response, err := nethttp.NewClient(nil).
		Send(nethttp.SetConf("http://localhost:7001/hello?name=christos", nethttp.Put), payload)

	replyOptions := network.NewOptions().Unmarshal(response.Options)

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
	options = network.NewOptions()
	options.SetHeader("Custom", "bar")
	payload = network.BuildPayload([]byte("foo"), options.Marshal())

	assert.NotPanics(t, func() { client.Send(netamqp.SetConf("/hello"), payload) },
		"Should not panics")

	// Create a payload
	options = network.NewOptions()

	options.SetHeader("Custom", "header-value")

	payload = network.BuildPayload([]byte("foo"), options.Marshal())

	// Send the payload
	response, err = netgrpc.NewClient("localhost:50051").
		Send(netgrpc.SetConf("/hello"), payload)

	replyOptions = network.NewOptions().Unmarshal(response.Options)

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
