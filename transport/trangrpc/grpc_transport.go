package trangrpc

import (
	"bitbucket.org/code_horse/pegasus/network"
	pb "bitbucket.org/code_horse/pegasus/transport/trangrpc/proto"
	"errors"
	"golang.org/x/net/context"
)

type Transporter struct {
	Router *Router
}

func (*Transporter) Send(properties *network.Properties, options *network.Options, payload []byte) (*network.Payload, error) {

	grpcProperties := NewProperties().BuildProperties(properties)

	connection := pb.NewServeClient(grpcProperties.GetConnection())

	grpcOptions := BuildOptions(options)

	if connection == nil {
		return nil, errors.New("Connection not found.")
	}

	syncResponse, err := connection.HandlerSync(context.Background(),
		&pb.HandlerRequest{
			Content: payload,
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

func (t *Transporter) Listen(properties *network.Properties, handler network.Handler, middleware network.Middleware) {
	grpcProperties := NewProperties().BuildProperties(properties)
	t.Router.Add(grpcProperties.GetPath(), handler, middleware)
}
