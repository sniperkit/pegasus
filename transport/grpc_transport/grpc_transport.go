package grpc_transport

import (
	"errors"
	pb "bitbucket.org/code_horse/pegasus/transport/grpc_transport/proto"
	"bitbucket.org/code_horse/pegasus/network"
	"golang.org/x/net/context"
)

type Transporter struct {
	Router *Router
}

func (*Transporter) Send(properties *network.Properties, options *network.Options, payload []byte) (*network.Payload, error) {

	grpcProperties := NewProperties().BuildProperties(properties)

	connection := grpcProperties.GetConnection()

	grpcOptions := BuildOptions(options)

	if connection == nil {
		return nil, errors.New("Connection not found.")
	}

	syncResponse, err := connection.HandlerSync(context.Background(),
		&pb.HandlerRequest{
			Content: []byte("hello world sync"),
			Options: grpcOptions.Marshal(),
			Path:    grpcProperties.GetPath(),
		},
	)

	if err != nil {
		return nil, err
	}

	pl := &network.Payload{Body: syncResponse.Content, Options: syncResponse.Options}

	return pl, nil
}

func (*Transporter) BuildStream(properties *network.Properties) (network.IStream, error) {

	grpcProperties := NewProperties().BuildProperties(properties)

	stream := &Stream{Path: grpcProperties.GetPath()}

	connection := grpcProperties.GetConnection()

	if connection == nil {
		return nil, errors.New("Connection not found. Set the network." +
			"Options.Private.GrpcConnection field.")
	}

	grpcStream, err := connection.Handler(context.Background())

	if err != nil {
		return nil, err
	}

	stream.Channel = grpcStream

	return stream, nil
}

func (t *Transporter) Listen(properties *network.Properties, handler network.Handler, middleware network.Middleware) {
	grpcProperties := NewProperties().BuildProperties(properties)
	t.Router.Add(grpcProperties.GetPath(), handler, middleware)
}
