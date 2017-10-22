package netamqp_test

import (
	"github.com/cpapidas/pegasus/network/netamqp"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

var NewConnection = netamqp.NewConnection

func TestNewConnection(t *testing.T) {
	var called bool

	// Set the mocked variable back to originals
	netamqp.Dial = amqp.Dial
	netamqp.NewConnection = NewConnection

	netamqp.RetriesTimes = 1
	netamqp.Sleep = 1

	called = false

	netamqp.Dial = func(url string) (*amqp.Connection, error) {
		called = true
		return &amqp.Connection{}, nil
	}

	connection, err := netamqp.NewConnection("")
	assert.NotNil(t, connection, "Connection should not be nil")
	assert.Nil(t, err, "error should be nil")
	assert.True(t, called, "Should call the Dial function")
}

func TestConnection_Channel(t *testing.T) {
	var called bool

	// Set the mocked variable back to originals
	netamqp.Dial = amqp.Dial
	netamqp.NewConnection = NewConnection

	netamqp.RetriesTimes = 1
	netamqp.Sleep = 1

	called = false

	netamqp.Dial = func(url string) (*amqp.Connection, error) {
		called = true
		c := &amqp.Connection{}
		return c, nil
	}

	connection, err := netamqp.NewConnection("")
	assert.Nil(t, err, "Error should be nil")
	assert.Panics(t, func() {
		connection.Channel()
	}, "Should panics for invalid connection")
}
