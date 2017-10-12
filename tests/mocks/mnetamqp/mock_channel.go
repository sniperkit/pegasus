package mnetamqp

import (
	"github.com/streadway/amqp"
)

// MockChannel mock for amqp.Channel
type MockChannel struct {
	CloseMock        func() error
	QueueDeclareMock func(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	PublishMock      func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	QosMock          func(prefetchCount, prefetchSize int, global bool) error
	ConsumeMock      func(
		queue,
		consumer string,
		autoAck,
		exclusive,
		noLocal,
		noWait bool,
		args amqp.Table,
	) (<-chan amqp.Delivery, error)
}

// Close mock for amqp.Channel Close function
func (m MockChannel) Close() error {
	if m.CloseMock != nil {
		return m.CloseMock()
	}
	return nil
}

// QueueDeclare mock for amqp.Channel QueueDeclare function
func (m MockChannel) QueueDeclare(
	name string,
	durable,
	autoDelete,
	exclusive,
	noWait bool,
	args amqp.Table,
) (amqp.Queue, error) {
	if m.QueueDeclareMock != nil {
		return m.QueueDeclareMock(name, durable, autoDelete, exclusive, noWait, args)
	}
	return amqp.Queue{}, nil
}

// Publish mock for amqp.Channel QueueDeclare function
func (m MockChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	if m.PublishMock != nil {
		return m.PublishMock(exchange, key, mandatory, immediate, msg)
	}
	return nil
}

// Qos mock for amqp.Channel Qos function
func (m MockChannel) Qos(prefetchCount, prefetchSize int, global bool) error {
	if m.QosMock != nil {
		return m.QosMock(prefetchCount, prefetchSize, global)
	}
	return nil
}

// Consume mock amqp.Channel Consume function
func (m MockChannel) Consume(
	queue,
	consumer string,
	autoAck,
	exclusive,
	noLocal,
	noWait bool,
	args amqp.Table,
) (<-chan amqp.Delivery, error) {
	if m.ConsumeMock != nil {
		return m.ConsumeMock(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	}
	return nil, nil
}
