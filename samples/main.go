package main

import (
	"fmt"
	"github.com/cpapidas/pegasus/samples/sample_grpc_http"
	"github.com/cpapidas/pegasus/samples/sample_grpc_http_amqp"
	"os"
)

func main() {

	if os.Args[1] == "sample_grpc_http" {
		switch os.Args[2] {
		case "client":
			sample_grpc_http.Client()
		case "server":
			sample_grpc_http.Server()
		default:
			fmt.Println("Command not found.")
		}
	} else if os.Args[1] == "sample_grpc_http_amqp" {
		switch os.Args[2] {
		case "client":
			sample_grpc_http_amqp.Client()
		case "server":
			sample_grpc_http_amqp.Server()
		default:
			fmt.Println("Command not found.")
		}
	} else {
		fmt.Println("Command not found.")
	}

}
