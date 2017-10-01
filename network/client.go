package network

// Client interface describes the protocols client models. It contains from two function the Send which is responsible
// to send the payload and the Close function which closes the opened connection is exists.
//
// Send Function
//
// Send function sends a payload to other servers.
//
// The first parameter is the conf which is a string array and could
// contains the URL or the method [POST, GET, ...]. In order to make it simpler each network provider (http, grpc, ...)
// have a SetPath(...) method which builds the conf parameter.
//
// The second parameter is the payload type of network.Payload. The payload could be created with function
// network.BuildPayload(). The payload.Options can be used in order to transport some options via options.Headers
// or some parameters via options.Params. The options.Body is similar with http body and used in order to
// transfer raw content.
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
// The Send function is will a payload object if server respond or an error if something go wrong.
//
// Close Function
//
// Close function closes the connection if exists and is open.
type Client interface {
	Send(conf []string, payload Payload) (*Payload, error)
	Close()
}
