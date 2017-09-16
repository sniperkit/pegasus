package dummies_rabitmq

import (
	"bitbucket.org/code_horse/pegasus/dummies"
	"bitbucket.org/code_horse/pegasus/helpers"
	"bitbucket.org/code_horse/pegasus/network"
	"bitbucket.org/code_horse/pegasus/transport"
	"bitbucket.org/code_horse/pegasus/transport/rabbitmq_transport"
	"fmt"
)

var RabbitMQ transport.ITransporter

func middleware(handler network.Handler, chanel *network.Channel) {
	handler(chanel)
}

func testHandler(channel *network.Channel) {
	payload := channel.Receive()
	options := rabbitmq_transport.BuildOptionsByBytes(payload.Options)

	output := "pegasus, RabitMQ Handler receive from /pegasus/sample/1 the data " + string(payload.Body) +
		" docker container: " + helpers.GetContainerId()
	fmt.Println(output, options)
	dummies.SendLogs(output)
}

func Init() {
	RabbitMQ = transport.ConnectToRabbitMQServer("amqp://guest:guest@haproxy:5672/")

	var properties *rabbitmq_transport.Properties

	properties = rabbitmq_transport.NewListenProperties().SetPath("/pegasus/sample/1").
		SetKey("/pegasus/sample/1")
	RabbitMQ.Listen(properties.GetProperties(), testHandler, middleware)
}
