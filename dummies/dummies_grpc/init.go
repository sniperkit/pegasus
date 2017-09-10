package dummies_grpc

import (
	"fmt"
	"bitbucket.org/code_horse/pegasus/dummies"
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/network"
	"bitbucket.org/code_horse/pegasus/transport"
	"bitbucket.org/code_horse/pegasus/transport/grpc_transport"
)

var grpcTransport transport.ITransporter

var grpcRouter transport.GrpcRouter

var demoServiceConnection transport.GrpcServeClient

func middleware(handler network.Handler, chanel *network.Channel) {
	fmt.Println("middleware...")
	handler(chanel)
}

func handlerSample1(channel *network.Channel) {

	params := channel.Receive()

	data := network.Payload{
		Body:    []byte("hello world. Got: " + string(params.Body)),
		Options: []byte(""),
	}

	output := "pegasus, GRPC Handler receive from /pegasus/sample/1 the data " + string(params.Body) +
		" docker container: " + helpers.GetContainerId()
	fmt.Println(output)
	dummies.SendLogs(output)

	channel.Send(data)
}

func Init() {
	grpcTransport, grpcRouter = transport.NewGrpcTransporter()

	properties := grpc_transport.NewProperties().SetPath("/pegasus/sample/1").GetProperties()

	grpcTransport.Listen(properties, handlerSample1, middleware)

	transport.StartGrpcServer(grpcRouter, "0.0.0.0:50051")
}

func ClientCall() {

	conn, err := transport.ConnectGrpcTo("localhost:50051")

	if err != nil {
		panic(err)
	}

	demoServiceConnection = transport.NewGrpcServeClient(conn)

	clientStreamCall()

	clientSyncConnection()
}

func clientStreamCall() {

	properties := grpc_transport.NewProperties().SetConnection(demoServiceConnection).SetPath("/pegasus/sample/1")

	stream, err := grpcTransport.BuildStream(properties.GetProperties())

	if err != nil {
		fmt.Println(err)
	}

	stream.Send(&network.Payload{Body: []byte("foo"), Options: grpc_transport.NewOptions().Marshal()})

	for {

		response, err := stream.Receive()

		if err != nil {
			panic(err)
		}

		if response.Body != nil {
			fmt.Println("Response from stream handler:", string(response.Body), string(response.Options))
			break
		}
	}

}

func clientSyncConnection() {

	properties := grpc_transport.NewProperties().SetPath("/pegasus/sample/1").SetConnection(demoServiceConnection)

	options := grpc_transport.NewOptions()

	pload, err := grpcTransport.Send(
		properties.GetProperties(),
		options.GetOptions(),
		[]byte("hello world from send method"),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Response from handler message:", string(pload.Body), string(pload.Options))
}
