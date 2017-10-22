package netamqp

import (
	"github.com/cpapidas/pegasus/peg"
	"github.com/streadway/amqp"
)

// RetriesTimes for reconnect with the server
var RetriesTimes = 10

// Sleep delay between calls
var Sleep = 5

// Server struct implements the peg.Server interface for RabbitMQ services.
type Server struct {
	connection IConnection
}

// NewServer returns a peg.Server object
var NewServer = func() peg.Server {
	return &Server{}
}

// SetConf gets a path as parameter and returns an array. It uses for Server.Listen.
func SetConf(path string) []string {
	return []string{path}
}

// Serve method (peg.Server) starts the RabbitMQ server for a specif address. It should have the right format
// <address>:<port>
func (s *Server) Serve(address string) {
	var err error

	peg.Retries(RetriesTimes, Sleep, func(...interface{}) bool {
		s.connection, err = NewConnection(address)
		if err != nil {
			return true
		}
		err = nil
		return false
	})

	if s.connection == nil || err != nil {
		panic("Cannot connect to RabbitMQ server")
	}
}

// Listen method starts a new worker which is listening to a specific queue.
func (s Server) Listen(conf []string, h peg.Handler, m peg.Middleware) {
	go s.addListener(conf[0], h, m)
}

// addListener add a new amqp listener
func (s Server) addListener(path string, handler peg.Handler, middleware peg.Middleware) {
	channel, err := s.connection.Channel()
	defer channel.Close()
	if err != nil {
		return
	}
	queue, err := s.queueDeclare(channel, path)
	if err != nil {
		return
	}
	err = channel.Qos(1, 0, false)
	if err != nil {
		return
	}
	msgs, err := s.createConsume(channel, queue.Name)
	if err != nil {
		return
	}
	forever := make(chan bool)
	go s.handlerRequest(msgs, handler, middleware)
	<-forever
}

// setHeader sets the amqp headers
func (Server) setHeader(options *peg.Options, headers amqp.Table) {
	for k, v := range headers {
		if peg.IsAMQPValidHeader(k) {
			options.SetHeader(k, v.(string))
		} else if paramKey := peg.AMQPParam(k); paramKey != "" {
			options.SetParam(paramKey, v.(string))
		}
	}
}

// handlerRequest is responsible to handle a specific request for the listener.
func (s Server) handlerRequest(msgs <-chan amqp.Delivery, handler peg.Handler, middleware peg.Middleware) {
	for d := range msgs {
		ch := peg.NewChannel(1)
		options := peg.NewOptions()
		s.setHeader(options, d.Headers)
		pl := peg.BuildPayload(d.Body, options.Marshal())
		ch.Send(pl)
		if middleware != nil {
			middleware(handler, ch)
		} else {
			handler(ch)
		}
		d.Ack(false)
	}
}

// queueDeclare declare a amqp queue and return the values.
func (Server) queueDeclare(channel IChannel, path string) (amqp.Queue, error) {
	return channel.QueueDeclare(
		path,
		true,
		false,
		false,
		false,
		nil,
	)
}

// createConsume create a new consumer and return the values.
func (Server) createConsume(channel IChannel, queueName string) (<-chan amqp.Delivery, error) {
	return channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}
