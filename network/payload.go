package network

// Payload struct define the payload that could be transferred over network protocol providers (HTTP, GRPC, AMQP)
//
// // Data tree:
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
type Payload struct {
	Body    []byte
	Options []byte
}

// NewPayload return a new payload object
func NewPayload(body []byte, options []byte) *Payload {
	return &Payload{Body: body, Options: options}
}

// BuildPayload return a new payload
func BuildPayload(body []byte, options []byte) Payload {
	return Payload{Body: body, Options: options}
}
