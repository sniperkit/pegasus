package transport

import (
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/transport/grpc_transport"
	pb "bitbucket.org/code_horse/pegasus/transport/grpc_transport/proto"
	"fmt"
	"google.golang.org/grpc"
)

type GrpcRouter *grpc_transport.Router

type GrpcServeClient pb.ServeClient

func NewGrpcTransporter() (ITransporter, *grpc_transport.Router) {
	router := &grpc_transport.Router{
		PathsWrapper: make(map[string]*grpc_transport.PathWrapper),
	}
	return &grpc_transport.Transporter{
		Router: router,
	}, router
}

func StartGrpcServer(router *grpc_transport.Router, address string) {
	server := &grpc_transport.Server{Router: router}
	server.Start(address)
}

func ConnectGrpcTo(path string) (*grpc.ClientConn, error) {

	var conn *grpc.ClientConn
	var err error

	fmt.Println("Start pool connect to GRPC")

	helpers.Retries(10, 5, func(...interface{}) bool {
		conn, err = grpc.Dial(path, grpc.WithInsecure())
		if err != nil {
			fmt.Println("Retring to connect to GRPC")
			return true
		}
		fmt.Println("Connected to GRPC at:", path)
		return false
	})

	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewGrpcServeClient(conn *grpc.ClientConn) pb.ServeClient {
	return pb.NewServeClient(conn)
}
