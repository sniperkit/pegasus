package rabbitmq_transport

import (
	"bitbucket.org/code_horse/pegasus/network"
	"github.com/streadway/amqp"
)

type Transport struct {
	Connection *amqp.Connection
}

func (t *Transport) Send(properties *network.Properties, options *network.Options, body []byte) (*network.Payload, error) {

	mqProperties := BuildRabbitMQProperties(properties)
	mqOptions := BuildOptions(options)

	consumer := NewComponent(mqProperties, t.Connection)
	ch, err := consumer.Connection.Channel()
	defer ch.Close()

	RabbitError(err, "Cannot connect")

	consumer.Channel = ch
	consumer.QueueDeclare()
	consumer.Publish(mqOptions, body)

	return nil, nil
}

func (t *Transport) Listen(properties *network.Properties, handler network.Handler, middleware network.Middleware) {
	propertiesRMQ := BuildRabbitMQProperties(properties)
	consumer := NewComponent(propertiesRMQ, t.Connection)
	go func() {
		ch, err := consumer.Connection.Channel()
		RabbitError(err, "Cannot connect")
		defer ch.Close()
		consumer.Channel = ch
		consumer.QueueDeclare()
		consumer.Qos()
		consumer.Consumer(handler, middleware)
	}()
}
