package netamqp_test

import (
	"github.com/cpapidas/pegasus/network"
	"github.com/cpapidas/pegasus/network/netamqp"
	"github.com/cpapidas/pegasus/tests/mocks/mnetamqp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/streadway/amqp"
	"time"
	"errors"
)

var _ = Describe("Client", func() {

	Describe("Client struct", func() {

		BeforeEach(func() {
			// Set the mocked variable back to originals
			netamqp.Dial = amqp.Dial
			netamqp.NewConnection = NewConnection
		})

		Context("Test constructor", func() {

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
				netamqp.NewClient("whatever")
				Expect(called).To(BeTrue())
			})

		})

		Context("Test Send function", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 1

			callChannel := false
			callClose := false
			callQueueDeclare := false
			callPublish := false

			mockConnection := &mnetamqp.MockConnection{}

			mockConnection.ChannelMock = func() (netamqp.IChannel, error) {
				mockChannel := &mnetamqp.MockChannel{}

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
						Expect(name).To(Equal("whatever/path"))
						Expect(durable).To(BeTrue())
						Expect(autoDelete).To(BeFalse())
						Expect(exclusive).To(BeFalse())
						Expect(noWait).To(BeFalse())
						Expect(args).To(BeNil())
					})

					return amqp.Queue{Name: "QueueName"}, nil
				}

				mockChannel.PublishMock = func(
					exchange,
					key string,
					mandatory,
					immediate bool,
					msg amqp.Publishing,
				) error {

					callPublish = true

					It("Should call the channel.Publish() with valid parameters", func() {

						Expect(exchange).To(BeEmpty())
						Expect(key).To(Equal("QueueName"))
						Expect(mandatory).To(BeFalse())
						Expect(immediate).To(BeFalse())
						Expect(immediate).To(BeFalse())

						Expect(msg.DeliveryMode).To(Equal(amqp.Persistent))
						Expect(msg.ContentType).To(Equal("text/plain"))

						Expect(msg.Headers["Sample"]).To(Equal("sample-content"))
						Expect(msg.Headers["HP-Sample"]).To(BeNil())
						Expect(msg.Headers["GR-Sample"]).To(BeNil())

						Expect(msg.Headers["MP-Foo"]).To(Equal("Bar"))

					})

					return nil
				}

				mockChannel.CloseMock = func() error {
					callClose = true
					return nil
				}

				return mockChannel, nil
			}

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return mockConnection, nil
			}

			client := netamqp.NewClient("whatever")

			options := network.NewOptions()
			options.SetHeader("Sample", "sample-content")
			options.SetHeader("HP-Sample", "sample-content")
			options.SetHeader("GR-Sample", "sample-content")

			options.SetParam("Foo", "Bar")

			payload := network.BuildPayload([]byte("body"), options.Marshal())

			client.Send(netamqp.SetConf("whatever/path"), payload)

			It("Should send the message", func() {
				Expect(callChannel).To(BeTrue())
			})

			It("Should call the close channel function", func() {
				Expect(callClose).To(BeTrue())
			})

			It("Should call the queue declare function", func() {
				Expect(callQueueDeclare).To(BeTrue())
			})

			It("Should call the publish function", func() {
				Expect(callPublish).To(BeTrue())
			})

		})

		Context("Test Send Configuration Publish", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 1

			callChannel := false
			callPublish := false

			mockConnection := &mnetamqp.MockConnection{}

			mockConnection.ChannelMock = func() (netamqp.IChannel, error) {
				mockChannel := &mnetamqp.MockChannel{}

				callChannel = true

				mockChannel.PublishMock = func(
					exchange,
					key string,
					mandatory,
					immediate bool,
					msg amqp.Publishing,
				) error {

					callPublish = true

					It("configure publishing method via headers", func() {
						Expect(msg.ContentType).To(Equal("application/json"))
						Expect(msg.ContentEncoding).To(Equal("content-encoding"))
						Expect(msg.DeliveryMode).To(Equal(uint8(1)))
						Expect(msg.Priority).To(Equal(uint8(2)))
						Expect(msg.CorrelationId).To(Equal("correlation-id"))
						Expect(msg.ReplyTo).To(Equal("reply-to"))
						Expect(msg.Expiration).To(Equal("expiration"))
						Expect(msg.MessageId).To(Equal("message-id"))
						t, _ := time.Parse("2006-01-02T15:04:05Z", "2006-01-02T15:04:05Z")
						Expect(msg.Timestamp).To(Equal(t))
						Expect(msg.Type).To(Equal("type"))
						Expect(msg.UserId).To(Equal("userid"))
						Expect(msg.AppId).To(Equal("appid"))
					})

					return nil
				}

				return mockChannel, nil
			}

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return mockConnection, nil
			}

			client := netamqp.NewClient("whatever")

			options := network.NewOptions()
			options.SetHeader("Content-Type", "application/json")
			options.SetHeader("MQ-Content-Encoding", "content-encoding")
			options.SetHeader("MQ-Delivery-Mode", "1")
			options.SetHeader("MQ-Priority", "2")
			options.SetHeader("MQ-Correlation-Id", "correlation-id")
			options.SetHeader("MQ-Reply-To", "reply-to")
			options.SetHeader("MQ-Expiration", "expiration")
			options.SetHeader("MQ-Message-Id", "message-id")
			options.SetHeader("MQ-Timestamp", "2006-01-02T15:04:05.000Z")
			options.SetHeader("MQ-Type", "type")
			options.SetHeader("MQ-User-Id", "userid")
			options.SetHeader("MQ-App-Id", "appid")

			payload := network.BuildPayload([]byte("body"), options.Marshal())

			client.Send(netamqp.SetConf("whatever/path"), payload)

			It("Should send the message", func() {
				Expect(callChannel).To(BeTrue())
			})

			It("Should call the channel.Publish function", func() {
				Expect(callPublish).To(BeTrue())
			})

		})

		Context("Test Close function", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 1

			callClose := false

			mockConnection := &mnetamqp.MockConnection{}

			mockConnection.CloseMock = func() error {
				callClose = true
				return nil
			}

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return mockConnection, nil
			}

			client := netamqp.NewClient("whatever")

			client.Close()

			It("Should call the Close function", func() {
				Expect(callClose).To(BeTrue())
			})

		})

		Context("Test on nil connection", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 1

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return nil, errors.New("error")
			}

			It("Should throw a panic on connection error", func() {
				Expect(func() {
					netamqp.NewClient("whatever")
				}).To(Panic())
			})

		})

		Context("Test Send function on channel failure", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 1

			mockConnection := &mnetamqp.MockConnection{}

			mockConnection.ChannelMock = func() (netamqp.IChannel, error) {
				return nil, errors.New("error")
			}

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return mockConnection, nil
			}

			client := netamqp.NewClient("whatever")

			payload := network.BuildPayload([]byte("body"), nil)

			_, err := client.Send(netamqp.SetConf("whatever/path"), payload)

			It("Should return an error", func() {
				Expect(err).ToNot(BeNil())
			})

		})

		Context("Test Send function on queue declaration failure", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 1

			callChannel := false
			callQueueDeclare := false

			mockConnection := &mnetamqp.MockConnection{}

			mockConnection.ChannelMock = func() (netamqp.IChannel, error) {
				mockChannel := &mnetamqp.MockChannel{}

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

				return mockChannel, nil
			}

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return mockConnection, nil
			}

			client := netamqp.NewClient("whatever")

			payload := network.BuildPayload([]byte("body"), nil)

			_, err := client.Send(netamqp.SetConf("whatever/path"), payload)

			It("Should call the channel function", func() {
				Expect(callChannel).ToNot(BeNil())
			})

			It("Should call the queue declaration function", func() {
				Expect(callQueueDeclare).ToNot(BeNil())
			})

			It("Should return an error", func() {
				Expect(err).ToNot(BeNil())
			})

		})

		Context("Test Send function on publish failure", func() {

			netamqp.RetriesTimes = 1
			netamqp.Sleep = 1

			callChannel := false
			callQueueDeclare := false
			callPublish := false

			mockConnection := &mnetamqp.MockConnection{}

			mockConnection.ChannelMock = func() (netamqp.IChannel, error) {
				mockChannel := &mnetamqp.MockChannel{}

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

				mockChannel.PublishMock = func(
					exchange,
					key string,
					mandatory,
					immediate bool,
					msg amqp.Publishing,
				) error {
					callPublish = true
					return errors.New("error")
				}

				return mockChannel, nil
			}

			netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
				return mockConnection, nil
			}

			client := netamqp.NewClient("whatever")

			payload := network.BuildPayload([]byte("body"), nil)

			_, err := client.Send(netamqp.SetConf("whatever/path"), payload)

			It("Should call the channel function", func() {
				Expect(callChannel).ToNot(BeNil())
			})

			It("Should call the queue declaration function", func() {
				Expect(callQueueDeclare).ToNot(BeNil())
			})

			It("Should call the publish function", func() {
				Expect(callPublish).ToNot(BeNil())
			})

			It("Should return an error", func() {
				Expect(err).ToNot(BeNil())
			})

		})

	})

})
