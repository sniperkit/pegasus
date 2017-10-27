package stress_test

import (
	"flag"
	"github.com/cpapidas/pegasus/netamqp"
	"github.com/cpapidas/pegasus/netgrpc"
	"github.com/cpapidas/pegasus/nethttp"
	"github.com/cpapidas/pegasus/peg"
	"testing"
	"time"
)

var (
	stress    = flag.Bool("stress", false, "run stress tests")
	amqp      = flag.Bool("amqp", false, "run stress tests include amqp")
	amqpBlock = make(chan bool)
)

const sleepInterval time.Duration = 2 * time.Second

func TestMain(m *testing.M) {
	flag.Parse()
	if !*stress {
		return
	}
	runServers(*amqp)
	time.Sleep(sleepInterval)
	m.Run()
}

func BenchmarkHTTP(b *testing.B) {
	client := nethttp.NewClient(nil)
	for n := 0; n < b.N; n++ {
		options := peg.NewOptions()
		options.SetHeader("Custom", "header-value")
		payload := peg.BuildPayload([]byte("foo"), options.Marshal())
		client.Send(nethttp.SetConf(":7001/hello?name=christos", nethttp.Put), payload)
	}
}

func BenchmarkGRPC(b *testing.B) {
	client := netgrpc.NewClient(":9001")
	for n := 0; n < b.N; n++ {
		options := peg.NewOptions()
		options.SetHeader("Custom", "header-value")
		payload := peg.BuildPayload([]byte("foo"), options.Marshal())
		client.Send(netgrpc.SetConf("/hello"), payload)
	}
}

func BenchmarkAMQP(b *testing.B) {
	if !*amqp {
		b.Skip()
	}
	client, _ := netamqp.NewClient("amqp://guest:guest@localhost:5672/")
	for n := 0; n < b.N; n++ {
		options := peg.NewOptions()
		options.SetHeader("Custom", "bar")
		payload := peg.BuildPayload([]byte("foo"), options.Marshal())
		client.Send(netamqp.SetConf("/hello"), payload)
		<- amqpBlock
	}
}

func runServers(rabbitMQ bool) {
	netamqp.RetriesTimes = 1
	netamqp.Sleep = 0
	serverAMQP := netamqp.NewServer()
	serverAMQP.Serve("amqp://guest:guest@localhost:5672/")

	serverHTTP := nethttp.NewServer(nil)
	serverGRPC := netgrpc.NewServer(nil)

	if rabbitMQ {
		serverAMQP = netamqp.NewServer()
		serverAMQP.Serve("amqp://guest:guest@localhost:5672/")
		serverAMQP.Listen(netamqp.SetConf("/hello"), amqpTestHandler, nil)
	}

	serverHTTP.Listen(nethttp.SetConf("/hello", nethttp.Put), generalTestHandler, nil)
	serverGRPC.Listen(netgrpc.SetConf("/hello"), generalTestHandler, nil)

	serverHTTP.Serve("localhost:7001")
	serverGRPC.Serve("localhost:9001")
}

func generalTestHandler(channel *peg.Channel) {
	// Receive the payload
	receive := channel.Receive()

	// Unmarshal options, change them and send them back
	options := peg.NewOptions().Unmarshal(receive.Options)

	replyOptions := peg.NewOptions()

	// RabbitMQ does not send back any response so we have to do the assertions inside handler
	if options.GetHeader("Custom") == "" || receive.Body == nil {
		panic("Header Custom and Body have to be set")
	}

	replyOptions.SetHeader("Custom", options.GetHeader("Custom")+" response")
	replyOptions.SetHeader("name", options.GetParam("name")+" response")
	replyOptions.SetHeader("id", options.GetParam("id")+" response")

	responseBody := string(receive.Body) + " response"

	// Create the new payload
	payload := peg.BuildPayload([]byte(responseBody), replyOptions.Marshal())

	// Send it back
	channel.Send(payload)
}

func amqpTestHandler(channel *peg.Channel) {
	generalTestHandler(channel)
	amqpBlock <- true
}
