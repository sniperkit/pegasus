package grpchttp

import (
	"fmt"
	"github.com/cpapidas/pegasus/network"
	"github.com/cpapidas/pegasus/network/netgrpc"
	"github.com/cpapidas/pegasus/network/nethttp"
)

func Client() {

	options := network.NewOptions()
	options.SetHeader("Token", "<th3S4cjgetTok42n>")
	options.SetParam("order", "DESC")

	send := network.BuildPayload([]byte("hello "), options.Marshal())

	// Send the http call and print the result
	httpResponse, err := nethttp.NewClient(nil).
		Send(nethttp.SetConf("http://localhost:9092/sample/14", nethttp.Put), send)
	if err != nil {
		panic(err)
	}

	// Send the grpc call and print the result
	grpcResponse, err := netgrpc.NewClient("localhost:9091").
		Send(netgrpc.SetConf("/sample/{id}"), send)

	if err != nil {
		panic(err)
	}

	// print the http response
	fmt.Println("--------------------- HTTP RESPONSE ---------------------")
	fmt.Println("body:", string(httpResponse.Body))
	fmt.Println("options:", network.NewOptions().Unmarshal(httpResponse.Options))
	fmt.Println("----------------------------------------------------------")

	// print the grpc response
	fmt.Println("--------------------- GRPC RESPONSE ---------------------")
	fmt.Println("body:", string(grpcResponse.Body))
	fmt.Println("options:", network.NewOptions().Unmarshal(grpcResponse.Options))
	fmt.Println("----------------------------------------------------------")
}
