package network

// Server is an interface that describes the Server struct. It's responsible to start the server and set up the
// listeners
type Server interface {

	// Serve function start the server for the configured router and giver address. The address must have the following
	// format <address>:<port>
	Serve(address string)

	// Listen function creates a handler for a specific endpoint. It gets the path string unique key, the handler
	// which is a function and the middleware which also is a function.
	//
	// Handler: 	(required) Handler is a function type of func(chanel *Channel). The Channel contains two method the
	// 				send method and the receive functions. The send method could send a payload to the other side
	// 				(server-client) and receive function return the payload.
	//
	//				Payload: The payload contains two fields:
	// 							Body: 		Which is the raw content
	//							Options: 	Which have the Header that tell to the server what to change and how to
	// 										process with the request and params which is set with all possible
	// 										parameters that can pass except the Body.
	// 										e.g. HTTP get prams, HTTP url params
	// 						Usually we build the payload.Options to network.Option content in order to have access to
	//						those fields.
	//						e.g. 	receive := channel.Receive() // return a network.Payload object
	//								options := network.NewOptions().Unmarshal(receive.Options)
	//								paramName := options.GetParam("name") // Get the name param
	//								options.SetHeader("foo", "bar") // Set a custom header
	//								p := network.BuildPayload([]byte("body"), options.Marshal())
	//								channel.Send(p)
	//
	// Middleware: 	 (optional) Middleware receive the network.Payload param and the network.Handler and executes first
	//				  This is responsible to add more parameter to handler and execute the handler function or not. At
	//				  Middleware for example you can check a authentication or authorization key, set the current user,
	//				  transform parameters, etc ...
	Listen(path []string, handler Handler, middleware Middleware)
}
