package nethttp_test

import (
	"github.com/cpapidas/pegasus/network"
	"github.com/cpapidas/pegasus/network/nethttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nethttp", func() {

	var handlerGet = func(channel *network.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := network.NewOptions().Unmarshal(receive.Options)

		if options.GetHeader("Content-Type") != "application/json" {
			panic("The header Content-Type should have default value application/json")
		}

		replyOptions := network.NewOptions()

		replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")

		// Create the new payload
		payload := network.BuildPayload([]byte(options.GetParam("foo")+" response"), replyOptions.Marshal())

		// Send it back
		channel.Send(payload)
	}

	var handlerPost = func(channel *network.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := network.NewOptions().Unmarshal(receive.Options)

		replyOptions := network.NewOptions()

		replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
		replyOptions.SetHeader("name", options.GetParam("name")+" response")

		responseBody := string(receive.Body) + " response"

		// Create the new payload
		payload := network.BuildPayload([]byte(responseBody), replyOptions.Marshal())

		// Send it back
		channel.Send(payload)
	}

	var handlerPut = func(channel *network.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := network.NewOptions().Unmarshal(receive.Options)

		replyOptions := network.NewOptions()

		replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
		replyOptions.SetHeader("name", options.GetParam("name")+" response")
		replyOptions.SetHeader("id", options.GetParam("id")+" response")

		responseBody := string(receive.Body) + " response"

		// Create the new payload
		payload := network.BuildPayload([]byte(responseBody), replyOptions.Marshal())

		// Send it back
		channel.Send(payload)
	}

	var handlerDelete = func(channel *network.Channel) {
		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := network.NewOptions().Unmarshal(receive.Options)

		replyOptions := network.NewOptions()

		replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
		replyOptions.SetHeader("name", options.GetParam("name")+" response")
		replyOptions.SetHeader("id", options.GetParam("id")+" response")

		// Create the new payload
		payload := network.BuildPayload([]byte(string(receive.Body)+" response"), replyOptions.Marshal())

		// Send it back
		channel.Send(payload)
	}

	var middleware = func(handler network.Handler, channel *network.Channel) {

		// Receive the payload
		receive := channel.Receive()

		// Unmarshal options, change them and send them back
		options := network.NewOptions().Unmarshal(receive.Options)

		options.SetHeader("Custom", options.GetHeader("Custom")+" middleware")

		// Create the new payload
		payload := network.BuildPayload(nil, options.Marshal())

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

	Describe("HTTP Server", func() {

		Context("Exchange message via HTTP", func() {

			It("Should not be nil", func() {
				Expect(server).ToNot(BeNil())
			})

		})

		Context("Send a GET request", func() {
			// Create a payload
			options := network.NewOptions()

			options.SetHeader("Custom", "header-value")
			options.SetHeader("MQ-Custom", "mq-value")
			options.SetHeader("GR-Custom", "gr-value")

			payload := network.BuildPayload(nil, options.Marshal())

			// Send the payload
			response, err := nethttp.NewClient(nil).
				Send(nethttp.SetConf("http://localhost:7000/http?foo=bar", nethttp.Get), payload)

			replyOptions := network.NewOptions().Unmarshal(response.Options)

			It("Should not throw an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return a application/json as default Content-Type", func() {
				Expect(replyOptions.GetHeader("Content-Type")).To(Equal("application/json"))
			})

			It("Should return nil headers for GR-* and MQ-*", func() {
				Expect(replyOptions.GetHeader("MQ-Custom")).To(BeEmpty())
				Expect(replyOptions.GetHeader("GR-Custom")).To(BeEmpty())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("bar response")))
				Expect(replyOptions.GetHeader("Custom")).To(Equal("header-value response"))
			})
		})

		Context("Send a POST request", func() {
			// Create a payload
			options := network.NewOptions()

			options.SetHeader("Custom", "header-value")

			payload := network.BuildPayload([]byte("foo"), options.Marshal())

			// Send the payload
			response, err := nethttp.NewClient(nil).
				Send(nethttp.SetConf("http://localhost:7000/http?name=christos", nethttp.Post), payload)

			replyOptions := network.NewOptions().Unmarshal(response.Options)

			It("Should not throw an error", func() {
				if err != nil {
					panic(err)
				}
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("foo response")))
				Expect(replyOptions.GetHeader("Custom")).To(Equal("header-value response"))
			})

			It("Should returns the param name", func() {
				Expect(replyOptions.GetHeader("Name")).To(Equal("christos response"))
			})
		})

		Context("Send a PUT request", func() {
			// Create a payload
			options := network.NewOptions()

			options.SetHeader("Custom", "header-value")

			payload := network.BuildPayload([]byte("foo"), options.Marshal())

			// Send the payload
			response, err := nethttp.NewClient(nil).
				Send(nethttp.SetConf("http://localhost:7000/http?name=christos", nethttp.Put), payload)

			replyOptions := network.NewOptions().Unmarshal(response.Options)

			It("Should not throw an error", func() {
				if err != nil {
					panic(err)
				}
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("foo response")))
				Expect(replyOptions.GetHeader("Custom")).To(Equal("header-value response"))
			})

			It("Should returns the param name", func() {
				Expect(replyOptions.GetHeader("Name")).To(Equal("christos response"))
			})
		})

		Context("Send a DELETE request", func() {
			// Create a payload
			options := network.NewOptions()

			options.SetHeader("Custom", "header-value")

			payload := network.BuildPayload([]byte("foo"), options.Marshal())

			// Send the payload
			response, err := nethttp.NewClient(nil).
				Send(nethttp.SetConf("http://localhost:7000/http?name=christos", nethttp.Delete), payload)

			replyOptions := network.NewOptions().Unmarshal(response.Options)

			It("Should not throw an error", func() {
				if err != nil {
					panic(err)
				}
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("foo response")))
				Expect(replyOptions.GetHeader("Custom")).To(Equal("header-value response"))
			})

			It("Should returns the param name", func() {
				Expect(replyOptions.GetHeader("Name")).To(Equal("christos response"))
			})

		})

		Context("Send a GET middleware request", func() {
			// Create a payload
			options := network.NewOptions()

			options.SetHeader("Custom", "header-value")

			payload := network.BuildPayload(nil, options.Marshal())

			// Send the payload
			response, err := nethttp.NewClient(nil).
				Send(nethttp.SetConf("http://localhost:7000/http/middleware?foo=bar", nethttp.Get), payload)

			It("Should not throw an error", func() {
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("bar response")))
				options := network.NewOptions().Unmarshal(response.Options)
				Expect(options.GetHeader("Custom")).To(Equal("header-value middleware response"))
			})
		})

	})

})
