package main

import (
	"fmt"
	"github.com/cpapidas/pegasus/examples/grpchttp"
	"github.com/cpapidas/pegasus/examples/grpchttpamqp"
	"os"
)

func main() {

	if os.Args[1] == "grpchttp" {
		switch os.Args[2] {
		case "client":
			grpchttp.Client()
		case "server":
			grpchttp.Server()
		default:
			fmt.Println("Command not found.")
		}
	} else if os.Args[1] == "grpchttpamqp" {
		switch os.Args[2] {
		case "client":
			grpchttpamqp.Client()
		case "server":
			grpchttpamqp.Server()
		default:
			fmt.Println("Command not found.")
		}
	} else {
		fmt.Println("Command not found.")
	}

}
