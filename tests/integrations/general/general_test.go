package general_test

import (
	"bitbucket.org/code_horse/pegasus/network"
	"bitbucket.org/code_horse/pegasus/network/netamqp"
	"bitbucket.org/code_horse/pegasus/network/netgrpc"
	"bitbucket.org/code_horse/pegasus/network/nethttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("General", func() {

	var skipAMQP = false

	Describe("Check rabbitMQ server", func() {
		defer func() {
			if r := recover(); r != nil {
				skipAMQP = true
			}
		}()

		netamqp.RetriesTimes = 1
		netamqp.Sleep = 0
		serverAMQP := netamqp.NewServer()
		serverAMQP.Serve("amqp://guest:guest@localhost:5672/")
	})

	Describe("Test all net* packages together", func() {

		if skipAMQP {
			PIt("RabbitMQ server is not listing")
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
		serverAMQP := netamqp.NewServer()

		serverAMQP.Serve("amqp://guest:guest@localhost:5672/")

		serverHTTP.Listen(nethttp.SetPath("/hello/{id}", nethttp.Put), handler, nil)
		serverGRPC.Listen(netgrpc.SetPath("/hello/{id}"), handler, nil)
		serverAMQP.Listen(netamqp.SetPath("/hello/{id}"), handler, nil)

		serverHTTP.Serve("localhost:7001")
		serverGRPC.Serve("localhost:50051")

		Context("Send a PUT request", func() {
			// Create a payload
			options := network.NewOptions()

			options.SetHeader("Custom", "header-value")

			payload := network.BuildPayload([]byte("foo"), options.Marshal())

			// Send the payload
			response, err := nethttp.NewClient().
				Send(nethttp.SetPath("http://localhost:7001/hello/44?name=christos", nethttp.Put), payload)

			replyOptions := network.NewOptions().Unmarshal(response.Options)

			It("Should not throw an error", func() {
				if err != nil {
					panic(err)
				}
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(string(response.Body)).To(Equal("foo response"))
				Expect(replyOptions.GetHeader("Custom")).To(Equal("header-value response"))
			})

			It("Should returns the param name", func() {
				Expect(replyOptions.GetHeader("Name")).To(Equal("christos response"))
			})

			It("Should return the path param", func() {
				Expect(replyOptions.GetHeader("Id")).To(Equal("44 response"))
			})
		})

		Context("Send AMQP request", func() {

			clientHTTP := netamqp.NewClient("amqp://guest:guest@localhost:5672/")

			It("Should not throw panic ", func() {
				options := network.NewOptions()
				options.SetHeader("Custom", "bar")
				payload := network.BuildPayload([]byte("foo"), options.Marshal())

				Expect(func() { clientHTTP.Send(netamqp.SetPath("/simple/handler"), payload) }).ToNot(Panic())
			})
		})

		Context("Send GRPC request", func() {

			// Create a payload
			options := network.NewOptions()

			options.SetHeader("Custom", "header-value")

			payload := network.BuildPayload([]byte("foo"), options.Marshal())

			// Send the payload
			response, err := netgrpc.NewClient("localhost:50051").
				Send(netgrpc.SetPath("/hello/{id}"), payload)

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
		})

	})

})
