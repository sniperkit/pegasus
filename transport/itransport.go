package transport

import (
	"bitbucket.org/code_horse/pegasus/network"
)

// ITransporter interface describes the abstract transportation methods Send and Receive. It's an interface for
// transport providers, each provider have to implement this interface in order to use it properly.
type ITransporter interface {

	// Send method sends the data to the client. The first params is the properties which are the private settings
	// for configuration, the seccond param is the options which are the public settings that client need to know
	// (Path, url params, etc, ...) and the third is the body which is the main payload. The body could be nil
	// and the properties and settings have at least configured.
	Send(properties *network.Properties, options *network.Options, payload []byte) (*network.Payload, error)

	// Listen method is the handler for a request. The first params is the properties which are the private settings and
	// configuration, the second param is handle function [func(chanel *Channel)] which gets a *network.Channel as
	// parameter in order to receive and send back the data and the middleware which is a common function which runs a
	// layer above. If the middleware is nil then the handler will execute only the handler function.
	Listen(properties *network.Properties, handler network.Handler, middleware network.Middleware)
}