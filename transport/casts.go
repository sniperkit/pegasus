package transport

import (
	"google.golang.org/grpc"
	"bitbucket.org/code_horse/pegasus/transport/trangrpc"
)

// ConnectionGRPC is the a type of GRPC connection object. It casted in order to use only the transport package.
type ConnectionGRPC *grpc.ClientConn

// Service is a type of identifier for a transportation service.
type Service string

// Router is an abstraction layer for transport routers
type Router interface{}

// Connection is the connection of transport
type Connection interface{}

type PropertiesGRPC trangrpc.Properties

const (
	// HTTP defines the http service
	HTTP Service = "HTTP"

	// GRPC defines the GRPC service
	GRPC Service = "GRPC"

	// RABBITMQ defines the RabbitMQ service
	RABBITMQ Service = "RABBITMQ"
)
