package rabbitmq_transport

import (
	"log"
	"fmt"
)

func RabbitError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
