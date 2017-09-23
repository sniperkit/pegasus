package tranrabbitmq

import (
	"fmt"
	"bitbucket.org/code_horse/pegasus/network"
	"github.com/streadway/amqp"
)

type Component struct {
	Properties *Properties
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

func NewComponent(properties *Properties, connection *amqp.Connection) *Component {
	return &Component{
		Properties: properties,
		Connection: connection,
	}
}

func (c *Component) ExchangeDeclare() {
	err := c.Channel.ExchangeDeclare(
		c.Properties.GetDeclarationName(),
		c.Properties.GetDeclarationType(),
		c.Properties.GetDeclarationDurable(),
		c.Properties.GetDeclarationAutoDelete(),
		c.Properties.GetDeclarationExclusive(),
		c.Properties.GetDeclarationNoWait(),
		nil,
	)
	RabbitError(err, "Cannot QueueDeclare")
}

func (c *Component) QueueDeclare() {
	queue, err := c.Channel.QueueDeclare(
		c.Properties.GetQueueName(),
		c.Properties.GetQueueDurable(),
		c.Properties.GetQueueAutoDelete(),
		c.Properties.GetQueueExclusive(),
		c.Properties.GetQueueNoWait(),
		nil,
	)
	fmt.Println("Listen QueueDeclare", c.Properties.GetQueueName())
	RabbitError(err, "Failed to declare a queue")
	c.Queue = queue
}

func (c *Component) QueueBind() {
	err := c.Channel.QueueBind(
		c.Queue.Name,
		c.Properties.GetKey(),
		c.Properties.GetDeclarationName(),
		c.Properties.GetQueueBindNoWait(),
		nil)
	RabbitError(err, "Failed to bind a queue")
}

func (c *Component) Qos() {
	err := c.Channel.Qos(
		c.Properties.GetQosPrefetchCount(),
		c.Properties.GetQosPrefetchSize(),
		c.Properties.GetQosGlobal(),
	)
	RabbitError(err, "Failed to set QoS")
}

func (c *Component) Consumer(handler network.Handler, middleware network.Middleware) {

	fmt.Println("Listen properties", c.Properties.GetConsumeName())

	msgs, err := c.Channel.Consume(
		c.Queue.Name,
		c.Properties.GetConsumeName(),
		c.Properties.GetConsumeAutoAct(),
		c.Properties.GetQueueExclusive(),
		c.Properties.GetConsumeNoLocal(),
		c.Properties.GetQueueNoWait(),
		nil,
	)
	RabbitError(err, "Cannot Consume")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			ch := network.NewChannel(2000)

			mqOptions := NewOptions().
				SetPublishContentType(d.ContentType).
				SetPublishDeliveryMode(d.DeliveryMode).
				SetPublishContentEncoding(d.ContentEncoding).
				SetPriority(d.Priority).
				SetPublishCorrelationId(d.CorrelationId).
				SetPublishReplyTo(d.ReplyTo).
				SetPublishExpiration(d.Expiration).
				SetPublishMessageId(d.MessageId).
				SetPublishTimestamp(d.Timestamp).
				SetPublishType(d.Type).
				SetPublishUserId(d.UserId).
				SetPublishAppId(d.AppId)

			pl := network.BuildPayload(d.Body, mqOptions.Marshal())
			ch.Send(pl)

			if middleware != nil {
				middleware(handler, ch)
			} else {
				handler(ch)
			}

			d.Ack(false)

		}

	}()
	<-forever
}

func (c *Component) Publish(options *Options, body []byte) {

	err := c.Channel.Publish(
		c.Properties.GetPublishExchange(),
		c.Queue.Name,
		c.Properties.GetPublishMandatory(),
		c.Properties.GetPublishImmediate(),
		amqp.Publishing{
			DeliveryMode:    options.GetPublishDeliveryMode(),
			ContentType:     options.GetPublishContentType(),
			ContentEncoding: options.GetPublishContentEncoding(),
			Priority:        options.GetPriority(),
			CorrelationId:   options.GetPublishCorrelationId(),
			ReplyTo:         options.GetPublishReplyTo(),
			Expiration:      options.GetPublishExpiration(),
			MessageId:       options.GetPublishMessageId(),
			Timestamp:       options.GetPublishTimestamp(),
			Type:            options.GetPublishContentType(),
			UserId:          options.GetPublishUserId(),
			AppId:           options.GetPublishAppId(),
			Body:            body,
		})

	RabbitError(err, "Cannot QueueDeclare")
}
