package transport

import (
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/transport/trangrpc"
	pb "bitbucket.org/code_horse/pegasus/transport/trangrpc/proto"
	"fmt"
	"google.golang.org/grpc"
)

type GrpcRouter *trangrpc.Router

type GrpcServeClient pb.ServeClient

func NewGrpcTransporter() (ITransporter, *trangrpc.Router) {
	router := &trangrpc.Router{
		PathsWrapper: make(map[string]*trangrpc.PathWrapper),
	}
	return &trangrpc.Transporter{
		Router: router,
	}, router
}

func StartGrpcServer(router *trangrpc.Router, address string) {
	server := &trangrpc.Server{Router: router}
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
