package peg

// Server is an interface that describes the Server implementation of transfer protocol providers. It's responsible to
// start the server and set up the listeners in order to process the requests.
//
// Serve
//
// Serve which is responsible to start the server. It gets an address string as parameters which should have the
// following format <address>:<port>
//
// Listen
//
// Listen function creates a handler for a specific endpoint. The parameters are
//
// conf: (required) Conf parameter describes the options that the listen needs in order to be created.
//
// e.g. The conf can be only a path but in case of HTTP protocol the method type is needed. Each protocol providers
// has a method called SetPath which is a conf.
//
// Handler: (required) Handler is a function type of func(chanel *peg.Channel). The Channel contains two functions
// the send method and the receive functions. The send function is used to send a payload and the receive function to
// receive the payload.
//
// Middleware: (optional) Middleware is a type of function which is executed before peg.handler function. It has
// two parameters the peg.Handler and the peg.channel. It is used only at peg.Server::Listen function.
// Usually the middleware type of function could be nil.
// The Middleware is responsible to call or not the handler function. Also it can edit the peg.Channel data that
// handler will get from the channel parameter.
//
// Data tree:
//  Payload:- contains two fields
// 		|- Body: Used to transfer the raw content.
//		|- Options: Options contains the Params, Headers, and Custom fields.
//			|
//			|- Params: We cannot set the params, only server can set this. Params could contain
//			| for example the url-path parameters.
//			|
//			|- Header: We can set headers in order to tell to server to change a behavior.
//			|
//			|- Fields: It contains custom fields for extreme cases.
//
// The payload also used as channel type through peg.Channel.
//
// Example:
//  var handlerGet = func(channel *peg.Channel) {
//  	// Receive the payload
//  	receive := channel.Receive()
//
//  	// Unmarshal options, change them and send them back
//  	options := peg.NewOptions().Unmarshal(receive.Options)
//
//  	replyOptions := peg.NewOptions()
//
//  	replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
//
//  	// Create the new payload
//  	payload := peg.BuildPayload([]byte(options.GetParam("foo")+" response"), replyOptions.Marshal())
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
//  options := peg.NewOptions()
//
//  options.SetHeader("Custom", "header-value")
//
//  payload := peg.BuildPayload(nil, options.Marshal())
//
//  // Send the payload
//  response, err := nethttp.NewClient().
//    Send(nethttp.SetPath("http://localhost:7000/http?foo=bar", nethttp.Get), payload)
type Server interface {
	Serve(address string)
	Listen(conf []string, handler Handler, middleware Middleware)
}
