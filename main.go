package main

import (
	"bitbucket.org/code_horse/pegasus/dummies/dummies_grpc"
	"bitbucket.org/code_horse/pegasus/dummies/dummies_http"
	"bitbucket.org/code_horse/pegasus/dummies/dummies_rabitmq"
)

var Forever chan bool

func main() {

	Forever = make(chan bool)

	//GRPC
	go dummies_grpc.Init()

	//RABBITMQ
	go dummies_rabitmq.Init()

	// HTTP
	go dummies_http.Init()

	<-Forever
}
