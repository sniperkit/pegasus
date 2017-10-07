package netamqp_test

import (
	"bitbucket.org/code_horse/pegasus/network"
	"bitbucket.org/code_horse/pegasus/network/netamqp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var failure = make(chan bool, 2)

var _ = Describe("Netamqp", func() {

	defer func() {
		if r := recover(); r != nil {
			PIt("RabbitMQ server is not listing")
		}
	}()

	var simpleHandler = func(channel *network.Channel) {

		receive := channel.Receive()

		options := network.BuildOptions(receive.Options)

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

	var middleware = func(handler network.Handler, channel *network.Channel) {
		receive := channel.Receive()

		options := network.BuildOptions(receive.Options)

		customHeaderValue := options.GetHeader("Custom") + " middleware"

		options.SetHeader("Custom", customHeaderValue)

		body := string(receive.Body) + " middleware"

		payload := network.BuildPayload([]byte(body), options.Marshal())

		channel.Send(payload)

		handler(channel)
	}

	netamqp.RetriesTimes = 1

	netamqp.Sleep = 0

	server := netamqp.NewServer()

	server.Serve("amqp://guest:guest@localhost:5672/")

	server.Listen(netamqp.SetConf("/simple/handler"), simpleHandler, middleware)

	client := netamqp.NewClient("amqp://guest:guest@localhost:5672/")

	Describe("RabbitMQ Integration Tests", func() {

		It("Should not throw panic ", func() {
			options := network.NewOptions()
			options.SetHeader("Custom", "bar")
			options.SetHeader("HP-Custom", "bar")
			options.SetHeader("GR-Custom", "bar")
			payload := network.BuildPayload([]byte("foo"), options.Marshal())

			Expect(func() { client.Send(netamqp.SetConf("/simple/handler"), payload) }).ToNot(Panic())

			Expect(<-failure).To(BeFalse())
		})

	})

})
