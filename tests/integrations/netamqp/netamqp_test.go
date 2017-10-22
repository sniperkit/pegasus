package netamqp_test

import (
	"github.com/cpapidas/pegasus/netamqp"
	"github.com/cpapidas/pegasus/peg"
	"github.com/stretchr/testify/assert"
	"testing"
)

var failure = make(chan bool, 2)

func TestNetamqp_integration(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skip(r)
		}
	}()

	var simpleHandler = func(channel *peg.Channel) {

		receive := channel.Receive()

		options := peg.BuildOptions(receive.Options)

		if options.GetHeader("HP-Custom") != "" || options.GetHeader("GR-Custom") != "" {
			panic("Header should contains key with prefix HP-* or GR-*")
		}

		if string(receive.Body) != "foo middleware" {
			failure <- true
			panic("body should be equal to foo middleware")
		}

		if options.GetHeader("Custom") != "bar middleware" {
			failure <- true
			panic("header should be equal to bae middleware")
		}

		failure <- false
	}

	var middleware = func(handler peg.Handler, channel *peg.Channel) {
		receive := channel.Receive()

		options := peg.BuildOptions(receive.Options)

		customHeaderValue := options.GetHeader("Custom") + " middleware"

		options.SetHeader("Custom", customHeaderValue)

		body := string(receive.Body) + " middleware"

		payload := peg.BuildPayload([]byte(body), options.Marshal())

		channel.Send(payload)

		handler(channel)
	}

	netamqp.RetriesTimes = 1

	netamqp.Sleep = 0

	server := netamqp.NewServer()

	server.Serve("amqp://guest:guest@localhost:5672/")

	server.Listen(netamqp.SetConf("/simple/handler"), simpleHandler, middleware)

	client, _ := netamqp.NewClient("amqp://guest:guest@localhost:5672/")

	// Should not throw panic
	options := peg.NewOptions()
	options.SetHeader("Custom", "bar")
	options.SetHeader("HP-Custom", "bar")
	options.SetHeader("GR-Custom", "bar")
	payload := peg.BuildPayload([]byte("foo"), options.Marshal())

	assert.NotPanics(t, func() { client.Send(netamqp.SetConf("/simple/handler"), payload) },
		"Should not panics")

	assert.False(t, <-failure, "Should be false")
}
