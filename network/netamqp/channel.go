package netamqp

import (
	"github.com/streadway/amqp"
)

// Channel describes the amqp.Channel. We use it in order to return a abstract object from connection.Channel().
type IChannel interface {
	Close() error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Qos(prefetchCount, prefetchSize int, global bool) error
	Consume(
		queue,
		consumer string,
		autoAck,
		exclusive,
		noLocal,
		noWait bool,
		args amqp.Table,
	) (<-chan amqp.Delivery, error)
}
