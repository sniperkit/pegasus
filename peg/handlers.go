package peg

// Handler is a type of function which has a parameter a peg.Channel. The handler was used at peg.Client::Send &
// peg.Server::Listen function in order to define a handler function which will process the request. Also handler
// is defined at peg.Middleware.
//
// Channel is the only and main way to transport data from Handlers to receivers. It contains a payload struct field
// with two nested fields payload.Options and payload.Body.
//
// Channel exists only in peg.Handler, peg.Middleware and peg.Client::Listen files, as parameter for those
// functions
//
// The payload.Options can be used in order to transport some options via options.Headers or some parameters via
// options.Params. The options.Body is similar with http body and is used in order to transfer raw content.
//
// An example for usage could be:
//	func handler(channel *peg.Channel) {
//		payload := channel.Receive() // Receive the payload here
//		options := peg.BuildOptions(payload.Options) // Convert the received options to options struct
// 		payload.Body = payload.Body + []byte(" sub-fix") // Change the body
//		chanel.send(peg.BuildPayload(payload.Body, options.Marshal())) // return the payload
//	}
type Handler func(chanel *Channel)

// Middleware is a type of function which is executed before peg.handler function. It has two parameters the
// peg.Handler and the peg.channel. It is used only at peg.Server::Listen function. Usually the middleware
// type of function would be nil.
//
// Handler parameter is used as handler function which will handler a specif request.
//
// Channel parameter is the only and main way to transport data from Handlers to receivers. It contains a payload struct
// field with two nested fields payload.Options and payload.Body.
//
// The Middleware is responsible to call or not the handler function. Also can edit the peg.Channel data that
// handler will get from channel parameter.
//
// An example for usage could be:
//	func middleware(handler Handler, chanel *Channel) {
//		payload := channel.Receive() // Receive the payload here
//		options := peg.BuildOptions(payload.Options) // Convert the received options to options struct
// 		options.SetHeader("Auth-Token", "@#$fsf2hkj42@#@#") // Add an extra header
// 		payload.Body = payload.Body + []byte(" sub-fix") // Change the body
//		chanel.send(peg.BuildPayload(payload.Body, options.Marshal())) // Add a new payload
// 		handler(chanel) // Call the handler with new channel
//	}
type Middleware func(handler Handler, channel *Channel)
