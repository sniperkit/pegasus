package netamqp

import (
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/network"
	"github.com/streadway/amqp"
)

// Client implements the network.Client interface. Client handles the remote calls via RabbitMQ (AMQP) protocol
type Client struct {
	connection *amqp.Connection
}

// NewClient connects to a RabbitMQ server, initializes a netamqp.Client object and returns a network.Client
var NewClient = func(address string) network.Client {

	var connection *amqp.Connection
	var err error

	helpers.Retries(RetriesTimes, Sleep, func(...interface{}) bool {
		connection, err = amqp.Dial(address)
		if err != nil {
			return true
		}
		return false
	})

	if connection == nil {
		// todo: [fix] [A002] Finish the Blunder package and throw an error
		panic("Cannot connect to RabbitMQ server")
	}

	return &Client{connection: connection}
}

// Send is responsible to send data in RabbitMQ server
func (c Client) Send(conf []string, payload network.Payload) (*network.Payload, error) {

	path := conf[0]
	body := payload.Body
	options := network.NewOptions().Unmarshal(payload.Options)

	// Create a new channel and make sure and channel will close when this function ends (defer)
	channel, err := c.connection.Channel()
	defer channel.Close()

	// todo: [fix] [A002] Finish the Blunder package and throw an error
	if err != nil {
		panic(err)
	}

	queue, err := channel.QueueDeclare(
		path,
		true,
		false,
		false,
		false,
		nil,
	)

	// todo: [fix] [A002] Finish the Blunder package and throw an error
	if err != nil {
		panic(err)
	}

	pub := publishing{}

	pub.DeliveryMode = amqp.Persistent
	pub.ContentType = "text/plain"

	pub.configurePublishing(options.GetHeaders())

	pub.Headers = make(map[string]interface{})

	// Set the headers
	for k, v := range options.GetHeaders() {
		pub.Headers[k] = v
	}

	pub.Body = body

	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		pub.Publishing,
	)

	// todo: [fix] [A002] Finish the Blunder package and throw an error
	if err != nil {
		panic(err)
	}

	return nil, nil
}

// Close close the current connection
func (c Client) Close() {
	c.connection.Close()
}
