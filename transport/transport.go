package transport

import (
	"bitbucket.org/code_horse/pegasus/network"
)

// This interface describe the abstract transportation methods that use in order to
// receive and send messages to third providers. Some of those transportation methods
// could be the http, grpc or even kafka and rabbitMQ.
type ITransporter interface {
	Send(properties *network.Properties, options *network.Options, payload []byte) (*network.Payload, error)

	BuildStream(properties *network.Properties) (network.IStream, error)

	Listen(properties *network.Properties, handler network.Handler, middleware network.Middleware)
}
