package netamqp_test

import (
	"errors"
	"github.com/cpapidas/pegasus/netamqp"
	"github.com/cpapidas/pegasus/peg"
	"github.com/cpapidas/pegasus/tests/mocks/mnetamqp"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

var serverNewConnection = netamqp.NewConnection

func setServerDefaults() {
	netamqp.Dial = amqp.Dial
	netamqp.NewConnection = serverNewConnection
}

func TestNewServer(t *testing.T) {
	setServerDefaults()
	// Should return a non empty object
	assert.NotNil(t, netamqp.NewServer(), "Should not return a nil server object")
}

func TestServer_Serve(t *testing.T) {
	setServerDefaults()
	// Should panics on connection not found
	netamqp.RetriesTimes = 0
	server := netamqp.NewServer()
	assert.Panics(t, func() {
		server.Serve("whatever")
	}, "Should panics on connection not found")

	// Should call the amqp.Dial method
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0
	called := false
	netamqp.Dial = func(url string) (*amqp.Connection, error) {
		called = true
		return &amqp.Connection{}, nil
	}
	server = netamqp.NewServer()
	server.Serve("")
	assert.True(t, called, "Should call the dial function")

	// Should panics on amqp.Dial error method
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0
	netamqp.Dial = func(url string) (*amqp.Connection, error) {
		called = true
		return nil, errors.New("error")
	}
	server = netamqp.NewServer()

	assert.Panics(t, func() { server.Serve("") }, "Should panic on Dial errors")
}

func TestServer_Listen(t *testing.T) {
	setClientDefaults()

	// Should call the amqp.Dial method
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callChannel := false
	callQueueDeclare := false
	callQos := false
	callConsume := false
	callHandler := false

	// Create a channel for mocks
	deliveries := make(chan amqp.Delivery)

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

		// Generate the mocked channel
		mockChannel := &mnetamqp.MockChannel{}

		callChannel = true

		// Mock the QueueDeclare function in order to check the parameters
		mockChannel.QueueDeclareMock = func(
			name string,
			durable,
			autoDelete,
			exclusive,
			noWait bool,
			args amqp.Table,
		) (amqp.Queue, error) {
			callQueueDeclare = true

			// Should has the right params
			assert.Equal(t, name, "/random/path")
			assert.True(t, durable, "Should durable be true")
			assert.False(t, autoDelete, "Should autoDelete be false")
			assert.False(t, exclusive, "Should exclusive be false")
			assert.False(t, noWait, "Should noWait be false")
			assert.Nil(t, args, "Should pass nil args")

			return amqp.Queue{Name: "QueueName"}, nil
		}

		// Mock the QosMock function in order to check the parameters
		mockChannel.QosMock = func(prefetchCount, prefetchSize int, global bool) error {

			callQos = true

			// Should has the right params
			assert.Equal(t, 1, prefetchCount, "Should be 1")
			assert.Equal(t, 0, prefetchSize, "Should be 0")
			assert.False(t, global, "Should be false")

			return nil
		}

		// Mock the Consume function in order to check the parameters
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

			// Should call the consume function with valid parameters
			assert.Equal(t, "QueueName", queue, "Should be equals to QueueName")
			assert.Empty(t, consumer, "Should be empty")
			assert.False(t, autoAck, "Should be false")
			assert.False(t, exclusive, "Should be false")
			assert.False(t, noLocal, "Should be false")
			assert.False(t, noWait, "Should be false")
			assert.Nil(t, args, "Should be nil")

			return (<-chan amqp.Delivery)(deliveries), nil
		}

		return mockChannel, nil
	}

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	server := netamqp.NewServer()
	server.Serve("whatever")

	c := make(chan bool)

	// The main handler
	var handler = func(channel *peg.Channel) {

		callHandler = true

		payload := channel.Receive()

		// Should return a body equals to body content
		assert.Equal(t, "body content", string(payload.Body),
			"Should be equals to body content")

		options := peg.NewOptions().Unmarshal(payload.Options)

		assert.Equal(t, "sample", options.GetHeader("Sample"), "Should be equals to Sample")
		assert.Empty(t, options.GetHeader("HP-Sample"), "Should be empty")
		assert.Empty(t, options.GetHeader("GR-Sample"), "Should be empty")
		assert.Equal(t, "sample", options.GetParam("Sample"), "Should be equals to sample")
		c <- true
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

	// wait unit handler was called
	<-c

	// Should call the Channel function=
	assert.True(t, callChannel, "Should call the function channel")

	// Should call the QueueDeclare function
	assert.True(t, callQueueDeclare, "Should call the function QueueDeclare")

	// Should call the Qos function"
	assert.True(t, callQos, "Should call the function Qos")

	// Should call the Consume function
	assert.True(t, callConsume, "Should call the function Consume")

	// Should call the handler
	assert.True(t, callHandler, "Should call the request Handler")
}

func TestServer_Listen_middleware(t *testing.T) {
	// Should call the amqp.Dial method
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callChannel := false
	callQueueDeclare := false
	callQos := false
	callConsume := false
	callHandler := false
	callMiddleware := false

	// Create a channel for mocks
	deliveries := make(chan amqp.Delivery)

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

		// Generate the mocked channel
		mockChannel := &mnetamqp.MockChannel{}

		callChannel = true

		// Mock the QueueDeclare function in order to check the parameters
		mockChannel.QueueDeclareMock = func(
			name string,
			durable,
			autoDelete,
			exclusive,
			noWait bool,
			args amqp.Table,
		) (amqp.Queue, error) {
			callQueueDeclare = true

			// Should has the right params
			assert.Equal(t, name, "/random/path")
			assert.True(t, durable, "Should durable be true")
			assert.False(t, autoDelete, "Should autoDelete be false")
			assert.False(t, exclusive, "Should exclusive be false")
			assert.False(t, noWait, "Should noWait be false")
			assert.Nil(t, args, "Should pass nil args")

			return amqp.Queue{Name: "QueueName"}, nil
		}

		// Mock the QosMock function in order to check the parameters
		mockChannel.QosMock = func(prefetchCount, prefetchSize int, global bool) error {

			callQos = true

			// Should has the right params
			assert.Equal(t, 1, prefetchCount, "Should be 1")
			assert.Equal(t, 0, prefetchSize, "Should be 0")
			assert.False(t, global, "Should be false")

			return nil
		}

		// Mock the Consume function in order to check the parameters
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

			// Should call the consume function with valid parameters
			assert.Equal(t, "QueueName", queue, "Should be equals to QueueName")
			assert.Empty(t, consumer, "Should be empty")
			assert.False(t, autoAck, "Should be false")
			assert.False(t, exclusive, "Should be false")
			assert.False(t, noLocal, "Should be false")
			assert.False(t, noWait, "Should be false")
			assert.Nil(t, args, "Should be nil")

			return (<-chan amqp.Delivery)(deliveries), nil
		}

		return mockChannel, nil
	}

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	server := netamqp.NewServer()
	server.Serve("whatever")

	c := make(chan bool)

	// The main handler
	var handler = func(channel *peg.Channel) {

		callHandler = true

		payload := channel.Receive()

		// Should return a body equals to body content
		assert.Equal(t, "body content", string(payload.Body),
			"Should be equals to body content")

		options := peg.NewOptions().Unmarshal(payload.Options)

		assert.Equal(t, "sample", options.GetHeader("Sample"), "Should be equals to Sample")
		assert.Empty(t, options.GetHeader("HP-Sample"), "Should be empty")
		assert.Empty(t, options.GetHeader("GR-Sample"), "Should be empty")
		assert.Equal(t, "middleware", options.GetHeader("middleware"),
			"Should be equals to middleware")
		assert.Equal(t, "sample", options.GetParam("Sample"), "Should be equals to sample")
		assert.Equal(t, "sample", options.GetParam("Sample"), "Should be equals to sample")
		assert.Equal(t, "middleware", options.GetParam("middleware"),
			"Should be equals to middleware")
		c <- true
	}

	// The middleware
	var middleware = func(handler peg.Handler, channel *peg.Channel) {
		callMiddleware = true
		receive := channel.Receive()

		options := peg.NewOptions().Unmarshal(receive.Options)
		options.SetParam("middleware", "middleware")
		options.SetHeader("middleware", "middleware")

		payload := peg.BuildPayload(receive.Body, options.Marshal())
		channel.Send(payload)
		handler(channel)
	}

	server.Listen(netamqp.SetConf("/random/path"), handler, middleware)

	headers := make(map[string]interface{})
	headers["Sample"] = "sample"
	headers["HP-Sample"] = "sample"
	headers["GR-Sample"] = "sample"
	headers["MP-Sample"] = "sample"

	deliveries <- amqp.Delivery{
		Body:    []byte("body content"),
		Headers: headers,
	}

	// wait unit handler was called
	<-c

	// Should call the Channel function=
	assert.True(t, callChannel, "Should call the function channel")

	// Should call the QueueDeclare function
	assert.True(t, callQueueDeclare, "Should call the function QueueDeclare")

	// Should call the Qos function"
	assert.True(t, callQos, "Should call the function Qos")

	// Should call the Consume function
	assert.True(t, callConsume, "Should call the function Consume")

	// Should call the middleware
	assert.True(t, callHandler, "Should call the middleware")

	// Should call the handler
	assert.True(t, callHandler, "Should call the request Handler")
}

func TestServer_Listen_channel(t *testing.T) {
	setServerDefaults()

	// Should call the amqp.Dial method
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callChannel := false
	callQueueDeclare := false

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	block := make(chan bool)

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

		// Generate the mocked channel
		mockChannel := &mnetamqp.MockChannel{}

		callChannel = true

		// Mock the QueueDeclare function in order to return an error
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
		block <- true
		return mockChannel, errors.New("error")
	}

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	server := netamqp.NewServer()
	server.Serve("whatever")

	server.Listen(netamqp.SetConf("/what/ever"), nil, nil)

	<-block

	// Should call the Channel function
	assert.True(t, callChannel, "Should call the channel function")

	// Should not call the QueueDeclare function
	assert.False(t, callQueueDeclare, "Should NOT call the QueueDeclare function")
}

func TestServer_Listen_queueDeclare(t *testing.T) {

	setServerDefaults()

	// Should call the amqp.Dial method
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callChannel := false
	callQueueDeclare := false
	callQos := false

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	block := make(chan bool)

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

		// Generate the mocked channel
		mockChannel := &mnetamqp.MockChannel{}

		callChannel = true

		// Mock the QueueDeclare function in order to return an error
		mockChannel.QueueDeclareMock = func(
			name string,
			durable,
			autoDelete,
			exclusive,
			noWait bool,
			args amqp.Table,
		) (amqp.Queue, error) {

			callQueueDeclare = true
			block <- true
			return amqp.Queue{}, errors.New("error")
		}

		mockChannel.QosMock = func(prefetchCount, prefetchSize int, global bool) error {
			callQos = true
			return nil
		}

		return mockChannel, nil
	}

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	server := netamqp.NewServer()
	server.Serve("whatever")

	server.Listen(netamqp.SetConf("/what/ever"), nil, nil)

	<-block

	// Should call the Channel function
	assert.True(t, callChannel, "Should call the channel function")

	// Should call the QueueDeclare function
	assert.True(t, callQueueDeclare, "Should call the QueueDeclare function")

	// Should not call the Qos function
	assert.False(t, callQos, "Should NOT call the Qos function")
}

func TestServer_Listen_qos(t *testing.T) {

	setServerDefaults()

	// Should call the amqp.Dial method
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callChannel := false
	callQueueDeclare := false
	callQos := false
	callConsume := false

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	block := make(chan bool)

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

		// Generate the mocked channel
		mockChannel := &mnetamqp.MockChannel{}

		callChannel = true

		// Mock the QueueDeclare function in order to return an error
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
			block <- true
			return errors.New("error")
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

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	server := netamqp.NewServer()
	server.Serve("whatever")

	server.Listen(netamqp.SetConf("/what/ever"), nil, nil)

	<-block

	// Should call the Channel function
	assert.True(t, callChannel, "Should call the channel function")

	// Should call the QueueDeclare function
	assert.True(t, callQueueDeclare, "Should call the QueueDeclare function")

	// Should call the Qos function
	assert.True(t, callQos, "Should call the Qos function")

	// Should NOT call the Consume function
	assert.False(t, callConsume, "Should NOT call the Consume function")
}

func TestServer_Listen_consume(t *testing.T) {

	setServerDefaults()

	// Should call the amqp.Dial method
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callChannel := false
	callQueueDeclare := false
	callQos := false
	callConsume := false

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	block := make(chan bool)

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

		// Generate the mocked channel
		mockChannel := &mnetamqp.MockChannel{}

		callChannel = true

		// Mock the QueueDeclare function in order to return an error
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
			block <- true
			return nil, errors.New("error")
		}

		return mockChannel, nil
	}

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	server := netamqp.NewServer()
	server.Serve("whatever")

	server.Listen(netamqp.SetConf("/what/ever"), nil, nil)

	<-block

	// Should call the Channel function
	assert.True(t, callChannel, "Should call the channel function")

	// Should call the QueueDeclare function
	assert.True(t, callQueueDeclare, "Should call the QueueDeclare function")

	// Should call the Qos function
	assert.True(t, callQos, "Should call the Qos function")

	// Should call the Consume function
	assert.True(t, callConsume, "Should call the Consume function")
}
