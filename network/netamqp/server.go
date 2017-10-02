package netamqp

import (
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/network"
	"github.com/streadway/amqp"
)

// RetriesTimes for reconnect with the server
var RetriesTimes = 10

// Sleep delay between calls
var Sleep = 5

// Server struct implements the network.Server interface for RabbitMQ services.
type Server struct {
	connection *amqp.Connection
}

// NewServer returns a network.Server object
var NewServer = func() network.Server {

	return &Server{}
}

// SetPath gets a path as parameter and returns an array. It uses for Server.Listen.
func SetPath(path string) []string {
	return []string{path}
}

// Serve method (network.Server) starts the RabbitMQ server for a specif address. It should have the right format
// <address>:<port>
func (s *Server) Serve(address string) {
	var connection *amqp.Connection
	var err error

	helpers.Retries(RetriesTimes, Sleep, func(...interface{}) bool {
		connection, err = amqp.Dial(address)
		if err != nil {
			return true
		}
		err = nil
		return false
	})

	if connection == nil || err != nil{
		// todo: [fix] [A002] Finish the Blunder package and throw an error
		panic("Cannot connect to RabbitMQ server")
	}

	s.connection = connection
}

// Listen method starts a new worker which is listening to a specific queue.
func (s Server) Listen(conf []string, handler network.Handler, middleware network.Middleware) {

	path := conf[0]

	go func() {

		if s.connection == nil {
			// todo: [fix] [A002] Finish the Blunder package and throw an error
			panic("RabbitMQ connection is nil, please start the server first and then set listeners")
		}

		// Create a channel
		channel, err := s.connection.Channel()

		// todo: [fix] [A002] Finish the Blunder package and throw an error
		if err != nil {
			panic(err)
		}
		defer channel.Close()

		queue, err := channel.QueueDeclare(
			path,
			true,
			false,
			false,
			false,
			nil,
		)

		err = channel.Qos(1, 0, false)

		// todo: [fix] [A002] Finish the Blunder package and throw an error
		if err != nil {
			panic(err)
		}

		msgs, err := channel.Consume(
			queue.Name,
			"",
			false,
			false,
			false,
			false,
			nil,
		)

		// todo: [fix] [A002] Finish the Blunder package and throw an error
		if err != nil {
			panic(err)
		}

		forever := make(chan bool)

		go func() {

			for d := range msgs {

				ch := network.NewChannel(1)

				options := network.NewOptions()

				// Set the headers
				for k, v := range d.Headers {
					options.SetHeader(k, v.(string))
				}

				pl := network.BuildPayload(d.Body, options.Marshal())
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

	}()
}
