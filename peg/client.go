package peg

// Client interface describes the protocols client models. It contains from two function the Send which is responsible
// to send the payload and the Close function which closes the opened connection is exists.
//
// Send Function
//
// Send function sends a payload to other servers.
//
// The first parameter is the conf which is a string array and could
// contain the URL or the method [POST, GET, ...]. In order to make it simpler each network provider (http, grpc, ...)
// has a SetPath(...) method which builds the conf parameter.
//
// The second parameter is the payload type of peg.Payload. The payload could be created by function
// peg.BuildPayload(). The payload.Options can be used in order to transport some options via options.Headers
// or some parameters via options.Params. The options.Body is similar with http body and is used in order to
// transfer raw content.
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
// The Send function will return a payload object. The server will respond with an error if something goes wrong.
//
// Close Function
//
// Close function will close the connection if it is open and exists.
type Client interface {
	Send(conf []string, payload Payload) (*Payload, error)
	Close() error
}
