package network

// Server is an interface that describes the Server implementation of transfer protocol providers. It's responsible to
// start the server and set up the listeners in order to process the requests.
//
// Serve
//
// Serve which is responsible to start the server. It gets a address string as parameters which should have the
// following format <address>:<port>
//
// Listen
//
// Listen function creates a handler for a specific endpoint. The parameters are
//
// conf: (required) Conf parameter describes the options that the listen needs in order to be created.
//
// e.g. The conf could be only a path but in case of HTTP protocol the method type is need. Each protocol providers
// has a method called SetPath which a conf.
//
// Handler: (required) Handler is a function type of func(chanel *network.Channel). The Channel contains two method the
// send method and the receive functions. The send function is used to send a payload and receive function to receive
// the payload.
//
// Middleware: (optional) Middleware is a type of function which executes before network.handler function. It has
// two parameters the network.Handler and the network.channel. It used only at network.Server::Listen function. Usually
// the middleware type of function could be nil.
// The Middleware is responsible to call or not the handler function. Also can edit the network.Channel data that
// handler will get from channel parameter.
//
// Data tree:
//  Payload:- contains two fields
// 		|- Body: Used to transfer the raw content.
//		|- Options: Options contains the Params, Headers a fields and Custom fields.
//			|
//			|- Params: We cannot set the params, only server can set this in order to return it. Params could contain
//			| for example the url-path parameters. If the param is set the server will ignore it.
//			|
//			|- Header: We can set headers in order to tell to server to change a behavior.
//			|
//			|- Fields: It contains custom fields for extreme cases.
//
// The payload also used as channel type through network.Channel.
//
// Example:
//  var handlerGet = func(channel *network.Channel) {
//  	// Receive the payload
//  	receive := channel.Receive()
//
//  	// Unmarshal options, change them and send them back
//  	options := network.NewOptions().Unmarshal(receive.Options)
//
//  	replyOptions := network.NewOptions()
//
//  	replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
//
//  	// Create the new payload
//  	payload := network.BuildPayload([]byte(options.GetParam("foo")+" response"), replyOptions.Marshal())
//
//  	// Send it back
//  	channel.Send(payload)
//  }
//
//  server := nethttp.NewServer(nil)
//  server.Listen(nethttp.SetPath("/http", nethttp.Get), handlerGet, nil)
//  server.Serve("localhost:7000")
//
//  // Send a request
//  // Create a payload
//  options := network.NewOptions()
//
//  options.SetHeader("Custom", "header-value")
//
//  payload := network.BuildPayload(nil, options.Marshal())
//
//  // Send the payload
//  response, err := nethttp.NewClient().
//    Send(nethttp.SetPath("http://localhost:7000/http?foo=bar", nethttp.Get), payload)
type Server interface {
	Serve(address string)
	Listen(conf []string, handler Handler, middleware Middleware)
}
