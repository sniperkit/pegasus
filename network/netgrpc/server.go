package netgrpc

import (
	"bitbucket.org/code_horse/pegasus/network"
	pb "bitbucket.org/code_horse/pegasus/transport/trangrpc/proto"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// Server is the manger of GRPC server. It's responsible to configure the server in order to run properly.
type Server struct {

	// Router handler the listener handlers
	Router *Router

	// address defines the server IP and port and has value with the following format (<IP>:<PORT>)
	address string

	// The server connection
	server *grpc.Server
}

// NewServer is the construct of Server struct. Generates and return a new Serve object. The parameter Router defines
// the route manager which is responsible to call the handlers.
var NewServer = func(router *Router) network.IServer {

	if router == nil {
		router = NewRouter()
	}

	return &Server{Router: router}
}

// SetPath get a path as parameter and return an array. It use for Server.Listen.
func SetPath(path string) []string {
	return []string{path}
}

// Serve function start a new GRPC server. It gets a string which is the address ("localhost:50099") and the GRPC router
// object which had been used before for listen methods that you want to configure.
func (s *Server) Serve(address string) {
	s.address = address
	go s.startServer()
}

// Listen function creates a handler for a specific endpoint. It gets the path string unique key, the handler
// which is a function and the middleware which also is a function.
func (s *Server) Listen(path []string, handler network.Handler, middleware network.Middleware) {
	s.Router.Add(path[0], handler, middleware)
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
// todo: [fix] [A005] Regenerate the protoc, and delete this method
func (s *Server) Handler(stream pb.Serve_HandlerServer) error {
	return errors.New("Not yet implement")
}

// startServer function start the server for the configured router and giver path
func (s *Server) startServer() {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.server = grpc.NewServer()

	pb.RegisterServeServer(s.server, &Server{Router: s.Router})

	// Register reflection service on gRPC server.
	reflection.Register(s.server)

	if err := s.server.Serve(lis); err != nil {
		// todo: [fix] [A002] Finish the Blunder package and throw an error
		log.Fatalf("failed to serve: %v", err)
	}
}
