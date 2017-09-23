package tranrabbitmq

import (
	"bitbucket.org/code_horse/pegasus/helpers"
	"fmt"
	"github.com/streadway/amqp"
)

// Server struct is responsible for RabbitMQ server. It manages connections and configuration may needed in order to
// ensure that the RabbitMQ server works properly
type Server struct {
}

// NewServer is constructor of Server struct. It initialize and return a Server object. It get a *mux.Router as
// parameter which could not be nil
func NewServer() *Server {
	return &Server{}
}

// Connect function connects to GRPC server
func (s *Server) Connect(path string) *amqp.Connection {
	var conn *amqp.Connection
	var err error

	// todo: [fix] [A004] Add those numbers to settings file
	helpers.Retries(10, 5, func(...interface{}) bool {
		conn, err = amqp.Dial(path)
		if err != nil {

			// todo: [fix] [A003] Fix the Blunder to print message if we want and add to it
			fmt.Println("Retring to connect to rabbitMQ")
			return true
		}
		return false
	})

	return conn
}
