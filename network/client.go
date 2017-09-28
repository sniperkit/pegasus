package network

// IClient interface describes the protocols client model. Client keep the connections open for each protocol in order
// close the connection should use the Close function. Send function used to send data to other servers.
type IClient interface {

	// Send function sends a payload to other servers. It gets the string path which is the unique id and the payload
	// object. The path may use other function ir order to generate the format for each provider.
	//
	// Payload: (required) contains two fields
	// 			Body: 		Here we can transfer the raw content of data
	// 			Options:	Options contains the Params and Headers a fields
	//
	// 				Params:	We cannot set the params only server is able to set this variable in order to return it to
	// 							us in order to handler it in Listen method. If the param is set the server will ignore
	// 							it.
	//
	// 				Header: We can set those header in order to tell the server to change a behavior. Server will return
	// 							to us the header after do the changes.
	//
	//						Header custom parameter: In header we have the Status key witch could get the HTTP Status.
	Send(path []string, payload Payload) (*Payload, error)

	// Close terminal the current connection.
	Close()
}
