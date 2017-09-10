package main

import (
	"bitbucket.org/code_horse/pegasus/dummies/dummies_http"
	"bitbucket.org/code_horse/pegasus/dummies/dummies_rabitmq"
	"bitbucket.org/code_horse/pegasus/dummies/dummies_grpc"
)

func main() {

	// GRPC
	go dummies_grpc.Init()

	// RABBITMQ
	dummies_rabitmq.Connect()
	dummies_rabitmq.Init()

	//// HTTP
	dummies_http.Init()
	dummies_http.StartServer()

	//fmt.Println("------ GRPC -------")
	//dummies_grpc.ClientCall()
	//fmt.Println("------ HTTP -------")
	//dummies_http.ClientCall()
	//fmt.Println("---- RABBITMQ -----")
	//dummies_rabitmq.ClientCall()
	//time.Sleep(time.Second * 4)
}
