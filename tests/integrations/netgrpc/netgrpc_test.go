package netgrpc_test

import (
	"bitbucket.org/code_horse/pegasus/network"
	"bitbucket.org/code_horse/pegasus/network/netgrpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("Netgrpc", func() {

	Describe("GRPC Server", func() {

		// Create a new server
		server := netgrpc.NewServer(nil)

		// Set the listeners

		server.Listen("/grpc/end-to-end", handler, nil)

		properties := netgrpc.NewProperties().SetPath("/grpc/end-to-end/properties")
		server.ListenTo(properties, handler, nil)

		server.Listen("/grpc/end-to-end/middleware", handler, middleware)

		properties = netgrpc.NewProperties().SetPath("/grpc/end-to-end/properties/middleware")
		server.ListenTo(properties, handler, middleware)

		// Start the server
		server.Serve("localhost:50052")

		Context("Exchange GRPC messages via Listen-Send hit at same time", func() {

			request := func(id string) {

				// Create a payload
				options := network.NewOptions()
				options.SetField("options", "value", id+"option")
				payload := network.BuildPayload([]byte(id+"hello message"), options.Marshal())

				// Send the payload
				netgrpc.NewClient("localhost:50052").Send("/grpc/end-to-end", payload)
				response, err := netgrpc.NewClient("localhost:50052").Send("/grpc/end-to-end", payload)

				It("Should not throw an error", func() {
					Expect(err).To(BeNil())
				})

				It("The response should have the following values", func() {
					Expect(response.Body).To(Equal([]byte(id + "hello message response")))
					options := network.NewOptions().Unmarshal(response.Options)
					Expect(options.Fields["options"]["value"]).To(Equal(id + "option response"))
				})
				fmt.Println(id)

			}

			go request("1")
			go request("2")

		})

		Context("Exchange GRPC payload via Listen-Send", func() {

			// Create a payload
			options := network.NewOptions()
			options.SetField("options", "value", "option")
			payload := network.BuildPayload([]byte("hello message"), options.Marshal())

			// Send the payload
			response, err := netgrpc.NewClient("localhost:50052").Send("/grpc/end-to-end", payload)

			It("Should not throw an error", func() {
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("hello message response")))
				options := network.NewOptions().Unmarshal(response.Options)
				Expect(options.Fields["options"]["value"]).To(Equal("option response"))
			})
		})

		Context("Exchange GRPC payload via ListenTo-Send", func() {
			// Create a payload
			options := network.NewOptions()
			options.SetField("options", "value", "option")
			payload := network.BuildPayload([]byte("hello message"), options.Marshal())

			// Send the payload
			response, err := netgrpc.NewClient("localhost:50052").
				Send("/grpc/end-to-end/properties", payload)

			It("Should not throw an error", func() {
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("hello message response")))
				options := network.NewOptions().Unmarshal(response.Options)
				Expect(options.Fields["options"]["value"]).To(Equal("option response"))
			})
		})

		Context("Exchange GRPC payload via Listen-Send with middleware", func() {
			// Create a payload
			options := network.NewOptions()
			options.SetField("options", "value", "option")
			payload := network.BuildPayload([]byte("hello message"), options.Marshal())

			// Send the payload
			response, err := netgrpc.NewClient("localhost:50052").
				Send("/grpc/end-to-end/middleware", payload)

			It("Should not throw an error", func() {
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("hello message middleware response")))
				options := network.NewOptions().Unmarshal(response.Options)
				Expect(options.Fields["options"]["value"]).To(Equal("option middleware response"))
			})
		})

		Context("Exchange GRPC payload via ListenTo-Send with middleware", func() {
			// Create a payload
			options := network.NewOptions()
			options.SetField("options", "value", "option")
			payload := network.BuildPayload([]byte("hello message"), options.Marshal())

			// Send the payload
			response, err := netgrpc.NewClient("localhost:50052").
				Send("/grpc/end-to-end/properties/middleware", payload)

			It("Should not throw an error", func() {
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("hello message middleware response")))
				options := network.NewOptions().Unmarshal(response.Options)
				Expect(options.Fields["options"]["value"]).To(Equal("option middleware response"))
			})
		})

		Context("Exchange GRPC payload via Listen-Send-Send hit the service twice", func() {

			// Create a payload
			options := network.NewOptions()
			options.SetField("options", "value", "option")
			payload := network.BuildPayload([]byte("hello message"), options.Marshal())

			// Send the payload
			netgrpc.NewClient("localhost:50052").Send("/grpc/end-to-end", payload)
			response, err := netgrpc.NewClient("localhost:50052").Send("/grpc/end-to-end", payload)

			It("Should not throw an error", func() {
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("hello message response")))
				options := network.NewOptions().Unmarshal(response.Options)
				Expect(options.Fields["options"]["value"]).To(Equal("option response"))
			})
		})

		Context("Exchange GRPC payload via Listen-Send with nil options in payload", func() {

			// Create a payload
			payload := network.BuildPayload([]byte("hello message"), nil)

			// Send the payload
			response, err := netgrpc.NewClient("localhost:50052").Send("/grpc/end-to-end", payload)

			It("Should not throw an error", func() {
				Expect(err).To(BeNil())
			})

			It("The response should have the following values", func() {
				Expect(response.Body).To(Equal([]byte("hello message response")))
			})
		})


	})

})

func handler(channel *network.Channel) {
	// Receive the payload
	receive := channel.Receive()

	// Unmarshal options, change them and send them back
	options := network.NewOptions().Unmarshal(receive.Options)
	options.SetField("options", "value", options.GetField("options", "value")+" response")

	// Create the new payload
	payload := network.BuildPayload([]byte(string(receive.Body)+" response"), options.Marshal())

	// Send it back
	channel.Send(payload)
}

func middleware(handler network.Handler, channel *network.Channel) {

	// Receive the payload
	receive := channel.Receive()

	// Unmarshal options, change them and send them back
	options := network.NewOptions().Unmarshal(receive.Options)
	options.SetField("options", "value", options.GetField("options", "value")+" middleware")

	// Create the new payload
	payload := network.BuildPayload([]byte(string(receive.Body)+" middleware"), options.Marshal())

	// Send it back
	channel.Send(payload)

	handler(channel)
}
