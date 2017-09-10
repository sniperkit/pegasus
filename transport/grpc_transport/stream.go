package grpc_transport

import (
	"bitbucket.org/code_horse/pegasus/network"
	pb "bitbucket.org/code_horse/pegasus/transport/grpc_transport/proto"
)

type Stream struct {
	Channel 	pb.Serve_HandlerClient
	Path string
}

func (s *Stream) Send(payload *network.Payload) {
	s.Channel.Send(&pb.HandlerRequest{Content: payload.Body, Options: payload.Options, Path: s.Path})
}

func (s *Stream) Receive() (*network.Payload, error) {
	receive, err :=  s.Channel.Recv()

	if err != nil{
		return nil, err
	}

	return &network.Payload{Body: receive.Content, Options: receive.Options}, nil
}
