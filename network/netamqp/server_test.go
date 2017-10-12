package netamqp_test

import (
	"github.com/cpapidas/pegasus/network/netamqp"

	"errors"
	"fmt"
	"github.com/cpapidas/pegasus/network"
	"github.com/cpapidas/pegasus/tests/mocks/mock_netamqp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/streadway/amqp"
)

var _ = Describe("Server", func() {

	Describe("Server constructor", func() {

		BeforeEach(func() {
			// Set the mocked variable back to originals
			netamqp.Dial = amqp.Dial
			netamqp.NewConnection = NewConnection
		})

		Context("Test struct constructor", func() {

			It("Should return a non empty object", func() {
				Expect(netamqp.NewServer()).ToNot(BeNil())
			})

		})

		Context("Test Serve function", func() {

			It("Should throw a panic when connection not found", func() {
				netamqp.RetriesTimes = 0
				Expect(func() { netamqp.NewClient("whatever") }).To(Panic())
			})

			It("Should call the amqp.Dial method", func() {
				netamqp.RetriesTimes = 1
				netamqp.Sleep = 1
				called := false
				netamqp.Dial = func(url string) (*amqp.Connection, error) {
					called = true
					return &amqp.Connection{}, nil
				}
				server := netamqp.NewServer()
				server.Serve("")
				Expect(called).To(BeTrue())
			})

		})

		Context("Test Listen method", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 0

			callChannel := false
			callQueueDeclare := false
			callQos := false
			callConsume := false
			callHandler := false

			deliveries := make(chan amqp.Delivery)

			mockConnection := &mock_netamqp.MockConnection{}

			mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

				mockChannel := &mock_netamqp.MockChannel{}

				callChannel = true

				mockChannel.QueueDeclareMock = func(
					name string,
					durable,
					autoDelete,
					exclusive,
					noWait bool,
					args amqp.Table,
				) (amqp.Queue, error) {

					callQueueDeclare = true

					It("Should call connection.Channel() with valid parameters", func() {
						Expect(name).To(Equal("/random/path"))
						Expect(durable).To(BeTrue())
						Expect(autoDelete).To(BeFalse())
						Expect(exclusive).To(BeFalse())
						Expect(noWait).To(BeFalse())
						Expect(args).To(BeNil())
					})

					return amqp.Queue{Name: "QueueName"}, nil
				}

				mockChannel.QosMock = func(prefetchCount, prefetchSize int, global bool) error {

					callQos = true

					It("Should call the Qos function with right parameters", func() {
						Expect(prefetchCount).To(Equal(1))
						Expect(prefetchSize).To(Equal(0))
						Expect(global).To(Equal(false))
					})

					return nil
				}

				mockChannel.ConsumeMock = func(
					queue,
					consumer string,
					autoAck,
					exclusive,
					noLocal,
					noWait bool,
					args amqp.Table,
				) (<-chan amqp.Delivery, error) {

					callConsume = true

					It("Should call the consume function with valid parameters", func() {
						Expect(queue).To(Equal("QueueName"))
						Expect(consumer).To(Equal(""))
						Expect(autoAck).To(BeFalse())
						Expect(exclusive).To(BeFalse())
						Expect(exclusive).To(BeFalse())
						Expect(noLocal).To(BeFalse())
						Expect(noWait).To(BeFalse())
						Expect(args).To(BeNil())
					})

					return (<-chan amqp.Delivery)(deliveries), nil
				}

				return mockChannel, nil
			}

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return mockConnection, nil
			}

			server := netamqp.NewServer()
			server.Serve("whatever")

			var handler = func(channel *network.Channel) {

				callHandler = true

				payload := channel.Receive()

				It("Should return a body equals to body content", func() {
					Expect(string(payload.Body)).To(Equal("body content"))
				})

				options := network.NewOptions().Unmarshal(payload.Options)

				It("Should have the right headers", func() {
					Expect(options.GetHeader("Sample")).To(Equal("sample"))
					Expect(options.GetHeader("HP-Sample")).To(Equal(""))
					Expect(options.GetHeader("GR-Sample")).To(Equal(""))
					fmt.Println("->>", options.GetParams())
					Expect(options.GetParam("Sample")).To(Equal("sample"))
				})

			}

			server.Listen(netamqp.SetConf("/random/path"), handler, nil)

			headers := make(map[string]interface{})
			headers["Sample"] = "sample"
			headers["HP-Sample"] = "sample"
			headers["GR-Sample"] = "sample"
			headers["MP-Sample"] = "sample"

			deliveries <- amqp.Delivery{
				Body:    []byte("body content"),
				Headers: headers,
			}

			It("Should call the Channel function", func() {
				Expect(callChannel).To(BeTrue())
			})

			It("Should call the QueueDeclare function", func() {
				Expect(callQueueDeclare).To(BeTrue())
			})

			It("Should call the Qos function", func() {
				Expect(callQos).To(BeTrue())
			})

			It("Should call the Consume function", func() {
				Expect(callConsume).To(BeTrue())
			})

			It("Should call the handler", func() {
				Expect(callHandler).To(BeTrue())
			})
		})

		Context("Test Channel throw an error", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 0

			callChannel := false
			callQueueDeclare := false

			mockConnection := &mock_netamqp.MockConnection{}

			mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

				mockChannel := &mock_netamqp.MockChannel{}

				callChannel = true

				mockChannel.QueueDeclareMock = func(
					name string,
					durable,
					autoDelete,
					exclusive,
					noWait bool,
					args amqp.Table,
				) (amqp.Queue, error) {

					callQueueDeclare = true

					return amqp.Queue{}, errors.New("error")
				}

				return mockChannel, errors.New("error")
			}

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return mockConnection, nil
			}

			server := netamqp.NewServer()
			server.Serve("whatever")

			server.Listen(netamqp.SetConf("/what/ever"), nil, nil)

			It("Should call the Channel function", func() {
				Expect(callChannel).To(BeTrue())
			})

			It("Should not call the QueueDeclare function", func() {
				Expect(callQueueDeclare).To(BeFalse())
			})

		})

		Context("Test Listen method error on QueueDeclare", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 0

			callChannel := false
			callQueueDeclare := false
			callQos := false

			mockConnection := &mock_netamqp.MockConnection{}

			mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

				mockChannel := &mock_netamqp.MockChannel{}

				callChannel = true

				mockChannel.QueueDeclareMock = func(
					name string,
					durable,
					autoDelete,
					exclusive,
					noWait bool,
					args amqp.Table,
				) (amqp.Queue, error) {

					callQueueDeclare = true

					return amqp.Queue{}, errors.New("error")
				}

				mockChannel.QosMock = func(prefetchCount, prefetchSize int, global bool) error {
					callQos = true
					return nil
				}

				return mockChannel, nil
			}

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return mockConnection, nil
			}

			server := netamqp.NewServer()
			server.Serve("whatever")

			server.Listen(netamqp.SetConf("/what/ever"), nil, nil)

			It("Should call the Channel function", func() {
				Expect(callChannel).To(BeTrue())
			})

			It("Should call the QueueDeclare function", func() {
				Expect(callQueueDeclare).To(BeTrue())
			})

			It("Should not clal the function Qos", func() {
				Expect(callQos).To(BeFalse())
			})

		})

		Context("Test Listen method error on Consume", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 0

			callChannel := false
			callQueueDeclare := false
			callQos := false
			callConsume := false

			mockConnection := &mock_netamqp.MockConnection{}

			mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

				mockChannel := &mock_netamqp.MockChannel{}

				callChannel = true

				mockChannel.QueueDeclareMock = func(
					name string,
					durable,
					autoDelete,
					exclusive,
					noWait bool,
					args amqp.Table,
				) (amqp.Queue, error) {

					callQueueDeclare = true

					return amqp.Queue{}, nil
				}

				mockChannel.QosMock = func(prefetchCount, prefetchSize int, global bool) error {
					callQos = true
					return nil
				}

				mockChannel.ConsumeMock = func(
					queue,
					consumer string,
					autoAck,
					exclusive,
					noLocal,
					noWait bool,
					args amqp.Table,
				) (<-chan amqp.Delivery, error) {
					callConsume = true
					return nil, errors.New("error")
				}

				return mockChannel, nil
			}

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return mockConnection, nil
			}

			server := netamqp.NewServer()
			server.Serve("whatever")

			server.Listen(netamqp.SetConf("/what/ever"), nil, nil)

			It("Should call the Channel function", func() {
				Expect(callChannel).To(BeTrue())
			})

			It("Should call the QueueDeclare function", func() {
				Expect(callQueueDeclare).To(BeTrue())
			})

			It("Should call the function Qos", func() {
				Expect(callQos).To(BeTrue())
			})

		})

	})

})
