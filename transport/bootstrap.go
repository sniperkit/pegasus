package transport

import (
	"bitbucket.org/code_horse/pegasus/transport/trangrpc"
	"bitbucket.org/code_horse/pegasus/transport/tranhttp"
	"bitbucket.org/code_horse/pegasus/transport/tranrabbitmq"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

// Bootstrap is the struct which is responsible to start a service. For example could start the HTTP or GRPC server
// Generally we should call the Bootstrap struct through the construct method.
type Bootstrap struct {
	// Service is a type of identifier for a transportation service. The hole struct will be configured according to
	// this property, so it cannot be nil.
	service Service

	// Router is an abstraction layer for transports
	router Router

	// Connection is the connection of transport
	connection Connection
}

// NewBootstrap creates a bootstrap object and returns it.
func NewBootstrap(service Service) *Bootstrap {
	return &Bootstrap{service: service}
}

// Serve finds the transport service and call the serve method. The params param is the path which could by anything
// according to service and the second param is the router which is an interface.
func (b *Bootstrap) Serve(path string) {

	switch b.service {
	case HTTP:
		tranhttp.NewServer(b.router.(*mux.Router)).Serve(path)
	case GRPC:
		go trangrpc.NewServer(b.router.(*trangrpc.Router)).Serve(path)
	default:
		// todo: [fix] [A002] Finish the Blunder package and throw an error
	}

}

// Generate and return a new transport interface.
func (b *Bootstrap) Generate() ITransporter {

	if b.service == HTTP {

		b.router = mux.NewRouter()

		return &tranhttp.Transporter{Router: b.router.(*mux.Router)}

	} else if b.service == GRPC {

		b.router = &trangrpc.Router{
			PathsWrapper: make(map[string]*trangrpc.PathWrapper),
		}

		return &trangrpc.Transporter{Router: b.router.(*trangrpc.Router)}

	} else if b.service == RABBITMQ {

		return &tranrabbitmq.Transport{Connection: b.connection.(*amqp.Connection)}

	}
	// todo: [fix] [A002] Finish the Blunder package and throw an error
	return nil
}

// SetRoute set the current Router
func (b *Bootstrap) SetRoute(router Router) *Bootstrap {
	b.router = router
	return b
}

// GetRoute return the current Router
func (b *Bootstrap) GetRoute() Router {
	return b
}

// Connect functions connects the current transport to another server.
func (b *Bootstrap) Connect(path string) Connection {

	switch b.service {
	case GRPC:
		b.connection = trangrpc.NewServer(b.router.(*trangrpc.Router)).Connect(path)
	case RABBITMQ:
		b.connection = tranrabbitmq.NewServer().Connect(path)
	default:
		// todo: [fix] [A002] Finish the Blunder package and throw an error
	}

	return b.connection
}
