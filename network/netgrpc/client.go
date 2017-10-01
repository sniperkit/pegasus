package netgrpc

import (
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/network"
	pb "bitbucket.org/code_horse/pegasus/network/netgrpc/proto"
	"context"
	"errors"
	"google.golang.org/grpc"
)

// Client implements the network.Client. Client struct describe the GRPC client. It contains all the functionality in
// order to talk to another server.
type Client struct {

	// Connection is the connection of GRPC server
	Connection *grpc.ClientConn
}

// NewClient connects to a GRPC server, set up the struct and return the new Client object
var NewClient = func(address string) network.Client {
	client := &Client{}
	client.Connection = client.connect(address)
	return client
}

// Send function sends a payload to a GRPC server. It gets the string path which is the unique id and the payload
// object.
func (c *Client) Send(path []string, payload network.Payload) (*network.Payload, error) {

	connection := pb.NewServeClient(c.Connection)

	if connection == nil {
		// todo: [fix] [A002] Finish the Blunder package and throw an error
		return nil, errors.New("CONNECTION NOT FOUND")
	}

	syncResponse, err := connection.HandlerSync(context.Background(),
		&pb.HandlerRequest{
			Content: payload.Body,
			Options: payload.Options,
			Path:    path[0],
		},
	)

	if err != nil {
		return nil, err
	}

	pl := &network.Payload{Body: syncResponse.Content, Options: syncResponse.Options}

	return pl, nil
}

// Close terminate the connection immediately.
func (c *Client) Close() {
	c.Connection.Close()
}

// connect used to connect with another GRPC service. The first parameter is address and returns the connection
// and ServeClient from proto buff
func (c *Client) connect(address string) *grpc.ClientConn {

	var conn *grpc.ClientConn
	var err error

	helpers.Retries(10, 5, func(...interface{}) bool {
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return true
		}
		return false
	})

	if err != nil {
		// todo: [fix] [A002] Finish the Blunder package and throw an error
		panic(err)
	}
	return conn
}
