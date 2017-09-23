package transport

import (
	"github.com/streadway/amqp"
	"bitbucket.org/code_horse/pegasus/helpers"
	"fmt"
	"bitbucket.org/code_horse/pegasus/transport/tranrabbitmq"
)

func ConnectToRabbitMQServer(path string) ITransporter{

	var conn *amqp.Connection
	var err error

	helpers.Retries(10, 5, func(...interface{}) bool {
		conn, err = amqp.Dial(path)
		if err != nil {
			fmt.Println("Retring to connect to rabbitMQ")
			return true
		}
		return false
	})

	tranrabbitmq.RabbitError(err, "Cannot connect to RabbitMQ server")
	return &tranrabbitmq.Transport{Connection: conn}
}