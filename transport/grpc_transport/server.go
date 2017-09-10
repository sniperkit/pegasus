package grpc_transport

import (
	pb "bitbucket.org/code_horse/pegasus/transport/grpc_transport/proto"
	"bitbucket.org/code_horse/pegasus/network"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
)

type Server struct {
	Router *Router
}

func (s *Server) Start(link string) {
	lis, err := net.Listen("tcp", link)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	pb.RegisterServeServer(server, &Server{Router: s.Router})
	// Register reflection service on gRPC server.
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (s *Server) HandlerSync(ctx context.Context, in *pb.HandlerRequest) (*pb.HandlerReply, error) {

	// Unmarshal the options
	options := &network.Options{}
	options.Unmarshal(in.Options)

	// Get the handler method for this request
	pathWrapper := s.Router.PathsWrapper[in.Path]

	params := network.Payload{Body: in.Content, Options: options.Marshal()}

	// Create a chanel
	channel := network.NewChannel(2000)

	channel.Send(params)

	// Check if path wrapper has a middleware function
	if pathWrapper.Middleware != nil {
		pathWrapper.Middleware(pathWrapper.Handler, channel)
	} else {
		pathWrapper.Handler(channel)
	}

	// Read from chanel
	payload := channel.Receive()

	// Return the content and the options to the client
	return &pb.HandlerReply{Content: payload.Body, Options: payload.Options}, nil
}

func (s *Server) Handler(stream pb.Serve_HandlerServer) error {

	// Create a chanel
	channel := network.NewChannel(2000)

	for {

		// Get the message
		in, err := stream.Recv()

		// If io exists
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		// Unmarshal the options
		options := &network.Options{}
		options.Unmarshal(in.Options)

		// Build the params object
		params := network.Payload{Body: in.Content, Options: options.Marshal()}

		// Set up a channel
		channel.Send(params)

		// Call the handler for given path
		pathWrapper := s.Router.PathsWrapper[in.Path]

		// Check if path wrapper has a middleware function
		if pathWrapper.Middleware != nil {
			pathWrapper.Middleware(pathWrapper.Handler, channel)
		} else {
			pathWrapper.Handler(channel)
		}

		// Read from chanel
		payload := channel.Receive()

		// Send the response from handler to the client
		stream.Send(&pb.HandlerReply{payload.Body, payload.Options})
	}
}
