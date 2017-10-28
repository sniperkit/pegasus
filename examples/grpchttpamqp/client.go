package grpchttpamqp

import (
	"fmt"

	"github.com/cpapidas/pegasus/netamqp"
	"github.com/cpapidas/pegasus/netgrpc"
	"github.com/cpapidas/pegasus/nethttp"
	"github.com/cpapidas/pegasus/peg"
)

// Client of grpc http amqp
func Client() {

	rabbitMQAddress := "amqp://guest:guest@localhost:5672/"

	options := peg.NewOptions()
	options.SetHeader("Token", "<th3S4cjgetTok42n>")
	options.SetParam("order", "DESC")

	send := peg.BuildPayload([]byte("hello "), options.Marshal())

	// Send the http call and print the result
	httpResponse, err := nethttp.NewClient("http://localhost:9092/").
		Send(nethttp.SetConf("sample", nethttp.Put), send)
	if err != nil {
		panic(err)
	}

	// Send the grpc call and print the result
	grpcResponse, err := netgrpc.NewClient("localhost:9091").
		Send(netgrpc.SetConf("/sample"), send)

	if err != nil {
		panic(err)
	}

	// Send the grpc call and print the result
	amqpClient, err := netamqp.NewClient(rabbitMQAddress)
	amqpClient.Send(netamqp.SetConf("/sample"), send)

	if err != nil {
		panic(err)
	}

	// print the http response
	fmt.Println("--------------------- HTTP RESPONSE ---------------------")
	fmt.Println("body:", string(httpResponse.Body))
	fmt.Println("options:", peg.NewOptions().Unmarshal(httpResponse.Options))
	fmt.Println("----------------------------------------------------------")

	// print the grpc response
	fmt.Println("--------------------- GRPC RESPONSE ---------------------")
	fmt.Println("body:", string(grpcResponse.Body))
	fmt.Println("options:", peg.NewOptions().Unmarshal(grpcResponse.Options))
	fmt.Println("----------------------------------------------------------")

	// print the grpc response
	fmt.Println("--------------------- GRPC RESPONSE ---------------------")
	fmt.Println("Go back to server terminal and check the logs")
	fmt.Println("----------------------------------------------------------")
}
