package dummies_http

import (
	"fmt"
	"bitbucket.org/code_horse/pegasus/transport"
	"bitbucket.org/code_horse/pegasus/transport/http_transport"
	"bitbucket.org/code_horse/pegasus/network"
	"github.com/gorilla/mux"
	"net/http"
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/dummies"
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
	// Create an HTTP router
	Router = mux.NewRouter()

	// Create a new HTTP transporter for our new router
	HttpTransporter = transport.NewHttpTransporter(Router)

	properties := http_transport.NewProperties().
		SetPath("/pegasus/sample/11").
		SetGetMethod()

	// Listen to the http router using HTTP transporter Listen Method
	HttpTransporter.Listen(properties.GetProperties(), handlerSample, middleware)

}

func StartServer() {
	err := http.ListenAndServe("0.0.0.0:8900", Router)

	if err != nil {
		panic(err)
	}
}

func ClientCall() {

	options := http_transport.NewOptions().SetHeader("foo", "bar")

	properties := http_transport.NewProperties().
		SetPath("http://localhost:8900/pegasus/sample/11?pa=11").
		SetGetMethod()

	payload, err := HttpTransporter.Send(properties.GetProperties(), options.GetOptions(), nil)

	if err != nil {
		panic(err)
	}

	bodyContent := payload.Body
	responseOptions := http_transport.NewOptions().BuildFromBytes(payload.Options)

	fmt.Println("http payload", string(bodyContent), responseOptions)
}
