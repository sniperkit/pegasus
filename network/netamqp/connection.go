package netamqp

import (
	"github.com/streadway/amqp"
)

// IConnection describe the amqp Connection
type IConnection interface {
	Close() error
	Channel() (IChannel, error)
}

// Dial is the function which create a connection object with RabbitMQ server
var Dial = amqp.Dial

// Connection struct implements the IConnection interface and describes the amqp.Connection struct. AMQP library does
// not contain interface in order to be mocked so we have to implement ours own hexagon layer for it.
type Connection struct {
	amqp.Connection
}

// NewConnection creates a Connection object, uses the amqp.Dial function in order to create a amqp.Connection object
// and embeds it into Connection struct. Finally it will return a IConnection object or an error.
var NewConnection = func(address string) (IConnection, error) {
	connection := &Connection{}
	dialConnection, err := Dial(address)
	if err != nil {
		return nil, err
	}
	connection.Connection = *dialConnection
	return connection, nil
}

// Channel return a new channel or an error.
func (c Connection) Channel() (IChannel, error) {
	return c.Connection.Channel()
}
