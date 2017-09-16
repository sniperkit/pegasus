package dummies_http

import (
	"bitbucket.org/code_horse/pegasus/dummies"
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/network"
	"bitbucket.org/code_horse/pegasus/transport"
	"bitbucket.org/code_horse/pegasus/transport/http_transport"
	"fmt"
	"github.com/gorilla/mux"
)

var HttpTransporter transport.ITransporter
var Router *mux.Router

func middleware(handler network.Handler, chanel *network.Channel) {
	fmt.Println("middleware...")
	handler(chanel)
}

func handlerSample(channel *network.Channel) {

	params := channel.Receive()

	httpOptions := http_transport.NewOptions().BuildFromBytes(params.Options)

	data := network.Payload{
		Body:    []byte("hello world. Got: " + string(params.Body) + " Header foo->" + httpOptions.GetHeaders()["Foo"] + " END."),
		Options: network.NewOptions().Marshal(),
	}

	output := "pegasus, HTTP Handler receive from /pegasus/sample/11 the data " + string(params.Body) +
		" docker container: " + helpers.GetContainerId()
	fmt.Println(output)
	dummies.SendLogs(output)

	channel.Send(data)
}

func Init() {
	Router = mux.NewRouter()

	HttpTransporter = transport.NewHttpTransporter(Router)

	properties := http_transport.NewProperties().
		SetPath("/pegasus/sample/11").
		SetGetMethod()

	HttpTransporter.Listen(properties.GetProperties(), handlerSample, middleware)

	transport.StartHTTP("0.0.0.0:8900", Router)
}
