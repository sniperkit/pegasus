package netamqp

import (
	"errors"
	"github.com/cpapidas/pegasus/peg"
	"github.com/streadway/amqp"
)

// Client implements the peg.Client interface. Client handles the remote calls via RabbitMQ (AMQP) protocol
type Client struct {
	connection IConnection
}

// NewClient connects to a RabbitMQ server, initializes a netamqp.Client object and returns a peg.Client
var NewClient = func(address string) (peg.Client, error) {

	var connection IConnection
	var err error

	peg.Retries(RetriesTimes, Sleep, func(...interface{}) bool {
		connection, err = NewConnection(address)
		if err != nil {
			return true
		}
		return false
	})

	if connection == nil {
		return nil, errors.New("cannot connect")
	}

	return &Client{connection: connection}, nil
}

// Send is responsible to send data in RabbitMQ server
func (c Client) Send(conf []string, payload peg.Payload) (*peg.Payload, error) {

	options := peg.NewOptions().Unmarshal(payload.Options)
	channel, err := c.connection.Channel()

	if err != nil {
		return nil, err
	}

	defer channel.Close()

	queue, err := c.queueDeclare(channel, conf[0])

	if err != nil {
		return nil, err
	}

	pub := buildPublishing(options, payload.Body)
	err = c.setPublish(channel, queue.Name, pub)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Close close the current connection
func (c Client) Close() error {
	return c.connection.Close()
}

// queueDeclare declares a new amqp Queue and returns it. If something goes wrong the function will return an error
// object.
func (Client) queueDeclare(channel IChannel, path string) (amqp.Queue, error) {
	return channel.QueueDeclare(
		path,
		true,
		false,
		false,
		false,
		nil,
	)
}

// setPublish setup the amqp publish function
func (Client) setPublish(channel IChannel, queueName string, publishing publishing) error {
	return channel.Publish(
		"",
		queueName,
		false,
		false,
		publishing.Publishing,
	)
}
