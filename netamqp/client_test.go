package netamqp_test

import (
	"errors"
	"github.com/cpapidas/pegasus/netamqp"
	"github.com/cpapidas/pegasus/peg"
	"github.com/cpapidas/pegasus/tests/mocks/mnetamqp"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var clientNewConnection = netamqp.NewConnection

// Sets the defaults
func setClientDefaults() {
	netamqp.Dial = amqp.Dial
	netamqp.NewConnection = clientNewConnection
}

func TestNewClient(t *testing.T) {
	setClientDefaults()

	// Should call the amqp.Dial method
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0
	called := false
	// Mock the Dial function in order to return a mocked connection object
	netamqp.Dial = func(url string) (*amqp.Connection, error) {
		called = true
		return &amqp.Connection{}, nil
	}
	netamqp.NewClient("whatever")
	assert.True(t, called, "Should call the Dial function")

	// Should panic on Dial error
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0
	// Mock the Dial function in order to return an error
	netamqp.Dial = func(url string) (*amqp.Connection, error) {
		return nil, errors.New("true")
	}
	client, err := netamqp.NewClient("whatever")
	assert.Nil(t, client, "Should return a nil Client on error")
	assert.NotNil(t, err, "Should return an error for undefined connection")
}

func TestClient_Send(t *testing.T) {
	setClientDefaults()

	// Set the peg.Retries properties
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callChannel := false
	callClose := false
	callQueueDeclare := false
	callPublish := false

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

		// Generate the mocked channel
		mockChannel := &mnetamqp.MockChannel{}

		callChannel = true

		// Mock the queueDeclare function, we have to check if the parameters had well configured
		mockChannel.QueueDeclareMock = func(
			name string,
			durable,
			autoDelete,
			exclusive,
			noWait bool,
			args amqp.Table,
		) (amqp.Queue, error) {

			callQueueDeclare = true

			// Should call connection.Channel() with valid parameters
			assert.Equal(t, "whatever/path", name, "Should have the name whatever/path")
			assert.True(t, durable, "Should be false")
			assert.False(t, autoDelete, "Should be false")
			assert.False(t, exclusive, "Should be false")
			assert.False(t, noWait, "Should be false")
			assert.Nil(t, args, "Should be nil")

			return amqp.Queue{Name: "QueueName"}, nil
		}

		// Mock the Publish function in order to check the parameters
		mockChannel.PublishMock = func(
			exchange,
			key string,
			mandatory,
			immediate bool,
			msg amqp.Publishing,
		) error {

			callPublish = true

			// Should call the channel.Publish() with valid parameters
			assert.Empty(t, exchange, "Should be empty")
			assert.Equal(t, "QueueName", key, "Should be equals to QueueName")
			assert.False(t, mandatory, "Should be false")
			assert.False(t, immediate, "Should be false")
			assert.Equal(t, amqp.Persistent, msg.DeliveryMode, "Should be equals to amqp.Persistent")
			assert.Equal(t, "text/plain", msg.ContentType, "Should be equals to text/plain")
			assert.Equal(t, "sample-content", msg.Headers["Sample"], "Should be equals to sample-content")
			assert.Nil(t, msg.Headers["HP-Sample"], "Should be nil for unsupported content")
			assert.Nil(t, msg.Headers["GR-Sample"], "Should be nil for unsupported content")
			assert.Equal(t, "Bar", msg.Headers["MP-Foo"], "Should be equals Bar")

			return nil
		}

		// Mock the Close function
		mockChannel.CloseMock = func() error {
			callClose = true
			return nil
		}

		return mockChannel, nil
	}

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	client, _ := netamqp.NewClient("whatever")
	options := peg.NewOptions()
	options.SetHeader("Sample", "sample-content")
	options.SetHeader("HP-Sample", "sample-content")
	options.SetHeader("GR-Sample", "sample-content")
	options.SetParam("Foo", "Bar")
	payload := peg.BuildPayload([]byte("body"), options.Marshal())

	client.Send(netamqp.SetConf("whatever/path"), payload)

	// Should send the message
	assert.True(t, callChannel, "Should call the function Channel")

	//"Should call the close channel function", func() {
	assert.True(t, callClose, "Should call the function Close")

	//"Should call the queue declare function", func() {
	assert.True(t, callQueueDeclare, "Should call the function QueueDeclare")

	//"Should call the publish function", func() {
	assert.True(t, callPublish, "Should call the function Publish")
}

func TestClient_Send_publishing(t *testing.T) {
	setClientDefaults()

	// Set the peg.Retries properties
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callChannel := false
	callPublish := false

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

		// Generate the mocked channel
		mockChannel := &mnetamqp.MockChannel{}

		callChannel = true

		// Mock the Publish function in order to check the parameters
		mockChannel.PublishMock = func(
			exchange,
			key string,
			mandatory,
			immediate bool,
			msg amqp.Publishing,
		) error {

			callPublish = true

			assert.Equal(t, "application/json", msg.ContentType,
				"Should be equals to application/json")
			assert.Equal(t, "content-encoding", msg.ContentEncoding,
				"Should be equals to content-encoding")
			assert.Equal(t, uint8(1), msg.DeliveryMode, "Should be equals to uint8(1)")
			assert.Equal(t, uint8(2), msg.Priority, "Should be equals to uint8(2)")
			assert.Equal(t, "correlation-id", msg.CorrelationId,
				"Should be equals to correlation-id")
			assert.Equal(t, "reply-to", msg.ReplyTo, "Should be equals to reply-to")
			assert.Equal(t, "expiration", msg.Expiration, "Should be equals to expiration")
			assert.Equal(t, "message-id", msg.MessageId, "Should be equals to message-id")
			tim, _ := time.Parse("2006-01-02T15:04:05Z", "2006-01-02T15:04:05Z")
			assert.Equal(t, tim, msg.Timestamp, "Should be equals to 2006-01-02T15:04:05Z")
			assert.Equal(t, "type", msg.Type, "Should be equals to type")
			assert.Equal(t, "userid", msg.UserId, "Should be equals to userid")
			assert.Equal(t, "appid", msg.AppId, "Should be equals to appid")

			return nil
		}

		return mockChannel, nil
	}

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	client, _ := netamqp.NewClient("whatever")

	options := peg.NewOptions()
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

	payload := peg.BuildPayload([]byte("body"), options.Marshal())

	client.Send(netamqp.SetConf("whatever/path"), payload)

	// Should send the message
	assert.True(t, callChannel, "Should call the Channel function")

	// Should call the channel.Publish function
	assert.True(t, callPublish, "Should call the Publish function")
}

func TestClient_Close(t *testing.T) {
	setClientDefaults()

	// Set the peg.Retries properties
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callClose := false

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	// Mock the close function
	mockConnection.CloseMock = func() error {
		callClose = true
		return nil
	}

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	client, _ := netamqp.NewClient("whatever")

	client.Close()

	// Should call the Close function
	assert.True(t, callClose, "Should call the Close function")
}

func TestClient_Send_channel(t *testing.T) {
	setClientDefaults()
	// Should configure the channel function

	// Set the peg.Retries properties
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {
		return nil, errors.New("error")
	}

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	client, _ := netamqp.NewClient("whatever")
	payload := peg.BuildPayload([]byte("body"), nil)
	_, err := client.Send(netamqp.SetConf("whatever/path"), payload)

	// Should return an error on channel failure
	assert.NotNil(t, err, "Should return an error on channel failure")
}

func TestClient_Send_queueDeclare(t *testing.T) {
	setClientDefaults()

	// Set the peg.Retries properties
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callChannel := false
	callQueueDeclare := false

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

		// Generate the mocked channel
		mockChannel := &mnetamqp.MockChannel{}

		callChannel = true

		// Mock the QueueDeclare function in order return an error
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

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	client, _ := netamqp.NewClient("whatever")

	payload := peg.BuildPayload([]byte("body"), nil)

	_, err := client.Send(netamqp.SetConf("whatever/path"), payload)

	// Should call the channel function
	assert.True(t, callChannel, "Should call the channel function")

	// Should call the queue declaration function
	assert.True(t, callQueueDeclare, "Should call the channel function")

	// Should return an error
	assert.NotNil(t, err, "Should return an error")
}

func TestClient_Send_publishFailure(t *testing.T) {
	setClientDefaults()

	// Set the peg.Retries properties
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0

	// Set the call* variables in order to test if the following methods are called
	callChannel := false
	callQueueDeclare := false
	callPublish := false

	// Generate a mocked connection object
	mockConnection := &mnetamqp.MockConnection{}

	// Have to configure the following functions in ChannelMock method
	mockConnection.ChannelMock = func() (netamqp.IChannel, error) {

		// Generate the mocked channel
		mockChannel := &mnetamqp.MockChannel{}

		callChannel = true

		// Mock the QueueDeclare function
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

		// Mock the Publish function in order to return an error
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

	// Replace the NewConnection function in order to run the mocked object
	netamqp.NewConnection = func(address string) (netamqp.IConnection, error) {
		return mockConnection, nil
	}

	client, _ := netamqp.NewClient("whatever")

	payload := peg.BuildPayload([]byte("body"), nil)

	_, err := client.Send(netamqp.SetConf("whatever/path"), payload)

	// Should call the channel function
	assert.True(t, callChannel, "Should call the channel function")

	// Should call the queue declaration function
	assert.True(t, callQueueDeclare, "Should call the channel function")

	// Should call the queue publish
	assert.True(t, callPublish, "Should call the publish function")

	// Should return an error
	assert.NotNil(t, err, "Should return an error")
}
