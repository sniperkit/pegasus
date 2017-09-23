package trangrpc

import (
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/network"
	pb "bitbucket.org/code_horse/pegasus/transport/trangrpc/proto"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type IServer interface {
	Serve(path string)
	Connect(path string) *grpc.ClientConn
}

// Server is the manger of GRPC server. It's responsible to configure the server in order to run properly.
type Server struct {

	// Router handler the listener handlers
	Router *Router
}

// NewServer is the construct of Server struct. Generates and return a new Serve object. The parameter Router defines
// the route manager which is responsible to call the handlers.
func NewServer(router *Router) IServer {

	if router == nil {
		router = &Router{}
	}

	return &Server{Router: router}
}

// Serve function start the server for the configured router and giver path
func (s *Server) Serve(path string) {
	lis, err := net.Listen("tcp", path)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	pb.RegisterServeServer(server, &Server{Router: s.Router})

	// Register reflection service on gRPC server.
	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		// todo: [fix] [A002] Finish the Blunder package and throw an error
		log.Fatalf("failed to serve: %v", err)
	}
}

// Connect used to connect with another GRPC service. The first parameter is path and returns the connection
// and ServeClient from proto buff
func (s *Server) Connect(path string) *grpc.ClientConn {

	var conn *grpc.ClientConn
	var err error

	helpers.Retries(10, 5, func(...interface{}) bool {
		conn, err = grpc.Dial(path, grpc.WithInsecure())
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

// Start Soon you will be deleted :)
// todo: [fix] [A003] delete it may not needed with new implementation of bootstrap
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

// HandlerSync is the method which send and receive the messages between the GRPC server. This is the brain of GRPC
// transportation. The first parameter is the context from GRPC the second is the HandlerRequest which container all
// the new messages. The function returns a HandlerReply or an error if something happens.
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

// Handler soon you will me deleted
// todo: [fix] [A003] delete it may not needed with new implementation of bootstrap
func (s *Server) Handler(stream pb.Serve_HandlerServer) error {
	return errors.New("Not yet implement")
}
