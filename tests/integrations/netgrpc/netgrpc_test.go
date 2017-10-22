package netgrpc_test

import (
	"github.com/cpapidas/pegasus/netgrpc"
	"github.com/cpapidas/pegasus/peg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetgrpc_integration(t *testing.T) {

	var handler = func(channel *peg.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := peg.NewOptions().Unmarshal(receive.Options)

		options.SetField("options", "value", options.GetField("options", "value")+" response")
		options.SetHeader("GR-Whatever", options.GetHeader("GR-Whatever"))

		// Create the new payload
		payload := peg.BuildPayload([]byte(string(receive.Body)+" response"), options.Marshal())

		// Send it back
		channel.Send(payload)
	}

	var middleware = func(handler peg.Handler, channel *peg.Channel) {

		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := peg.NewOptions().Unmarshal(receive.Options)
		options.SetField("options", "value", options.GetField("options", "value")+" middleware")

		// Create the new payload
		payload := peg.BuildPayload([]byte(string(receive.Body)+" middleware"), options.Marshal())

		// Send it back
		channel.Send(payload)

		handler(channel)
	}

	// Create a new server
	server := netgrpc.NewServer(nil)

	// Set the listeners

	server.Listen(netgrpc.SetConf("/grpc/end-to-end"), handler, nil)

	server.Listen(netgrpc.SetConf("/grpc/end-to-end/middleware"), handler, middleware)

	// Start the server
	server.Serve("localhost:50052")

	// Exchange GRPC messages via Listen-Send hit at same time

	request := func(id string) {

		// Create a payload
		options := peg.NewOptions()
		options.SetField("options", "value", id+"option")
		payload := peg.BuildPayload([]byte(id+"hello message"), options.Marshal())

		// Send the payload
		netgrpc.NewClient("localhost:50052").Send([]string{"/grpc/end-to-end"}, payload)
		response, err := netgrpc.NewClient("localhost:50052").Send([]string{"/grpc/end-to-end"}, payload)

		// Should not throw an error
		assert.Nil(t, err, "Should be nil")

		// The response should have the following values
		assert.Equal(t, []byte(id+"hello message response"), response.Body,
			"Should be equals to id + hello message response")
		options = peg.NewOptions().Unmarshal(response.Options)
		assert.Equal(t, id+"option response", options.Fields["options"]["value"],
			"Should be equals to id + option response")
	}

	go request("1")
	go request("2")

	// Exchange GRPC payload via Listen-Send

	// Create a payload
	options := peg.NewOptions()
	options.SetField("options", "value", "option")
	options.SetHeader("HP-Whatever", "MQ-value")
	options.SetHeader("MQ-Whatever", "HP-value")
	options.SetHeader("GR-Whatever", "GR-value")
	payload := peg.BuildPayload([]byte("hello message"), options.Marshal())

	// Send the payload
	response, err := netgrpc.NewClient("localhost:50052").Send([]string{"/grpc/end-to-end"}, payload)

	replyOptions := peg.NewOptions().Unmarshal(response.Options)

	// Should not throw an error
	assert.Nil(t, err, "Should be nil")

	// Should return nil headers for HP-* and MQ-*
	assert.Empty(t, replyOptions.GetHeader("MQ-Custom"), "Should be empty")
	assert.Empty(t, replyOptions.GetHeader("HP-Custom"), "Should be empty")

	// Should return GR-* Header
	assert.Equal(t, "GR-value", replyOptions.GetHeader("GR-Whatever"),
		"Should be equals to GR-value")

	// The response should have the following values
	assert.Equal(t, []byte("hello message response"), response.Body,
		"Should be equals to hello message response")

	assert.Equal(t, "option response", replyOptions.Fields["options"]["value"],
		"Should be equals to option response")

	// Exchange GRPC payload via Listen-Send with middleware

	// Create a payload
	options = peg.NewOptions()
	options.SetField("options", "value", "option")
	payload = peg.BuildPayload([]byte("hello message"), options.Marshal())

	// Send the payload
	response, err = netgrpc.NewClient("localhost:50052").
		Send([]string{"/grpc/end-to-end/middleware"}, payload)

	// Should not throw an error
	assert.Nil(t, err, "Should be nil")

	// The response should have the following values
	assert.Equal(t, []byte("hello message middleware response"), response.Body,
		"Should be equals to hello message middleware response")
	options = peg.NewOptions().Unmarshal(response.Options)
	assert.Equal(t, "option middleware response", options.Fields["options"]["value"],
		"Should be equals to option middleware response")

	// Exchange GRPC payload via Listen-Send-Send hit the service twice

	// Create a payload
	options = peg.NewOptions()
	options.SetField("options", "value", "option")
	payload = peg.BuildPayload([]byte("hello message"), options.Marshal())

	// Send the payload
	netgrpc.NewClient("localhost:50052").Send([]string{"/grpc/end-to-end"}, payload)
	response, err = netgrpc.NewClient("localhost:50052").Send([]string{"/grpc/end-to-end"}, payload)

	// Should not throw an error
	assert.Nil(t, err, "Should be nil")

	// The response should have the following values
	assert.Equal(t, []byte("hello message response"), response.Body,
		"Should be hello message response")
	options = peg.NewOptions().Unmarshal(response.Options)
	assert.Equal(t, "option response", options.Fields["options"]["value"],
		"Should be equals to option response")

	// Exchange GRPC payload via Listen-Send with nil options in payload

	// Create a payload
	payload = peg.BuildPayload([]byte("hello message"), nil)

	// Send the payload
	response, err = netgrpc.NewClient("localhost:50052").Send([]string{"/grpc/end-to-end"}, payload)

	// Should not throw an error
	assert.Nil(t, err, "Should be nil")

	// The response should have the following values
	assert.Equal(t, []byte("hello message response"), response.Body,
		"Should be equals to hello message response")
}
