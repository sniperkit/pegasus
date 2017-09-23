package dummies_grpc

import (
	"bitbucket.org/code_horse/pegasus/dummies"
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/network"
	"bitbucket.org/code_horse/pegasus/transport"
	"bitbucket.org/code_horse/pegasus/transport/trangrpc"
	"fmt"
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
		" docker container: " + helpers.GetContainerID()
	fmt.Println(output)
	dummies.SendLogs(output)

	channel.Send(data)
}

func Init() {
	grpcTransport, grpcRouter = transport.NewGrpcTransporter()

	properties := trangrpc.NewProperties().SetPath("/pegasus/sample/1").GetProperties()

	grpcTransport.Listen(properties, handlerSample1, middleware)

	transport.StartGrpcServer(grpcRouter, "0.0.0.0:50051")

}
