package dummies_rabitmq

import (
	"fmt"
	"bitbucket.org/code_horse/pegasus/dummies"
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/network"
	"bitbucket.org/code_horse/pegasus/transport"
	"bitbucket.org/code_horse/pegasus/transport/rabbitmq_transport"
)

var RabbitMQ transport.ITransporter

func Connect() {
	RabbitMQ = transport.ConnectToRabbitMQServer("amqp://guest:guest@haproxy:5672/")
}

func middleware(handler network.Handler, chanel *network.Channel) {
	handler(chanel)
}

func testHandler(channel *network.Channel) {
	payload := channel.Receive()
	options := rabbitmq_transport.BuildOptionsByBytes(payload.Options)

	output := "pegasus, RabitMQ Handler receive from /pegasus/sample/11 the data " + string(payload.Body) +
		" docker container: " + helpers.GetContainerId()
	fmt.Println(output, options)
	dummies.SendLogs(output)
}

func Init() {
	properties := rabbitmq_transport.NewListenProperties().SetPath("task_queue").SetKey("task_queue")
	RabbitMQ.Listen(properties.GetProperties(), testHandler, middleware)
}

func ClientCall() {
	properties := rabbitmq_transport.NewSendProperties().SetPath("task_queue")
	options := rabbitmq_transport.NewSendOptions()
	RabbitMQ.Send(properties.GetProperties(), options.GetOptions(), []byte("Hello from RabbitMQ"))
}
