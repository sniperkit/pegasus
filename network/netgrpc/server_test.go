package netgrpc_test

import (
	"github.com/cpapidas/pegasus/network/netgrpc"

	"errors"
	"fmt"
	"github.com/cpapidas/pegasus/network"
	pb "github.com/cpapidas/pegasus/network/netgrpc/proto"
	"google.golang.org/grpc"
	"net"
	"reflect"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer(t *testing.T) {
	server := netgrpc.NewServer(nil)
	// Should return a new server instance of *netgrpc.Server
	assert.NotNil(t, server, "Should return a new server instance of *netgrpc.Server")
	assert.Equal(t, "<*netgrpc.Server Value>", reflect.ValueOf(server).String(),
		"Should be equals to <*netgrpc.Server Value>")
}

func TestSetConf(t *testing.T) {
	path := netgrpc.SetConf("path")
	// Should return the path as array
	assert.Equal(t, []string{"path"}, path, "Should be equals to ['path']")
}

func TestServer_Listen(t *testing.T) {

	var handler = func(channel *network.Channel) {}
	var middleware = func(handler network.Handler, channel *network.Channel) {}

	router := netgrpc.NewRouter()

	server := netgrpc.NewServer(router)
	server.Listen(netgrpc.SetConf("whatever"), handler, middleware)

	// Should add a new route
	assert.NotNil(t, router.PathsWrapper["whatever"], "Should not be nil")
}

func TestServer_HandlerSync(t *testing.T) {

	// Define the call* variable in order to make sure that methods have been called
	callHandler := false

	options := network.NewOptions()
	options.SetHeader("GR-Sample", "sample")
	options.SetHeader("HP-Sample", "sample")
	options.SetHeader("MQ-Sample", "sample")
	options.SetHeader("Custom-Sample", "sample")

	// Create an HandlerRequest object
	handlerRequest := &pb.HandlerRequest{
		Content: []byte("content"),
		Options: options.Marshal(),
		Path:    "the/path",
	}

	// The main handler
	var handler = func(channel *network.Channel) {

		callHandler = true

		payload := channel.Receive()
		options := network.NewOptions().Unmarshal(payload.Options)

		// Should receive valid params
		assert.Equal(t, "sample", options.GetHeader("Custom-Sample"),
			"Should be equals to sample")
		assert.Equal(t, "sample", options.GetHeader("GR-Sample"),
			"Should be equals to sample")
		assert.Empty(t, options.GetHeader("HP-Sample"),
			"Should be empty (invalid header)")
		assert.Empty(t, options.GetHeader("MQ-Sample"),
			"Should be empty (invalid header)")
		assert.Equal(t, "content", string(payload.Body),
			"Should be equals to content")

		replyOptions := network.NewOptions()
		replyOptions.SetHeader("Custom-Sample-Reply", "sample reply")
		replyOptions.SetHeader("GR-Sample-Reply", "sample reply")
		replyOptions.SetHeader("HP-Sample-Reply", "sample reply")
		replyOptions.SetHeader("MQ-Sample-Reply", "sample reply")

		channel.Send(network.BuildPayload([]byte("content-reply"), replyOptions.Marshal()))
	}

	router := netgrpc.NewRouter()

	server := &netgrpc.Server{Router: router}
	server.Listen(netgrpc.SetConf("the/path"), handler, nil)

	// Should add a new route
	assert.NotNil(t, router.PathsWrapper["the/path"], "Should not be nil")

	handlerReply, err := server.HandlerSync(nil, handlerRequest)

	// Should call the handler function
	assert.True(t, callHandler, "Should be true")

	// Should return nil error object
	assert.Nil(t, err, "Should be nil")

	// Should return the right values

	optionsUn := network.NewOptions().Unmarshal(handlerReply.Options)

	assert.Equal(t, "sample reply", optionsUn.GetHeader("Custom-Sample-Reply"),
		"Should be equals to sample reply")
	assert.Equal(t, "sample reply", optionsUn.GetHeader("GR-Sample-Reply"),
		"Should be equals to sample reply")
	assert.Empty(t, optionsUn.GetHeader("HP-Sample-Reply"),
		"Should be empty")
	assert.Empty(t, optionsUn.GetHeader("MQ-Sample-Reply"),
		"Should be empty")
	assert.Equal(t, "content-reply", string(handlerReply.Content),
		"Should be equal to content-reply")
}

func TestServer_HandlerSync_middleware(t *testing.T) {

	// Define the call* variable in order to make sure that methods have been called
	callHandler := false

	options := network.NewOptions()
	options.SetHeader("GR-Sample", "sample")
	options.SetHeader("HP-Sample", "sample")
	options.SetHeader("MQ-Sample", "sample")
	options.SetHeader("Custom-Sample", "sample")

	// Create an HandlerRequest object
	handlerRequest := &pb.HandlerRequest{
		Content: []byte("content"),
		Options: options.Marshal(),
		Path:    "the/path",
	}

	// The main handler
	var handler = func(channel *network.Channel) {

		callHandler = true

		payload := channel.Receive()
		options := network.NewOptions().Unmarshal(payload.Options)

		// Should receive valid params
		assert.Equal(t, "sample reply middleware", options.GetHeader("Custom-Sample"),
			"Should be equals to sample")
		assert.Equal(t, "sample reply middleware", options.GetHeader("GR-Sample"),
			"Should be equals to sample")

		replyOptions := network.NewOptions()
		replyOptions.SetHeader("Custom-Sample-Reply", "sample reply")
		replyOptions.SetHeader("GR-Sample-Reply", "sample reply")
		replyOptions.SetHeader("HP-Sample-Reply", "sample reply")
		replyOptions.SetHeader("MQ-Sample-Reply", "sample reply")

		channel.Send(network.BuildPayload([]byte("content-reply"), replyOptions.Marshal()))
	}

	var middleware = func(handler network.Handler, channel *network.Channel) {

		callHandler = true

		payload := channel.Receive()
		options := network.NewOptions().Unmarshal(payload.Options)

		// Should receive valid params
		assert.Equal(t, "sample", options.GetHeader("Custom-Sample"),
			"Should be equals to sample")
		assert.Equal(t, "sample", options.GetHeader("GR-Sample"),
			"Should be equals to sample")
		assert.Empty(t, options.GetHeader("HP-Sample"),
			"Should be empty (invalid header)")
		assert.Empty(t, options.GetHeader("MQ-Sample"),
			"Should be empty (invalid header)")
		assert.Equal(t, "content", string(payload.Body),
			"Should be equals to content")

		replyOptions := network.NewOptions()
		replyOptions.SetHeader("Custom-Sample", "sample reply middleware")
		replyOptions.SetHeader("GR-Sample", "sample reply middleware")

		channel.Send(network.BuildPayload([]byte("content"), replyOptions.Marshal()))
		handler(channel)
	}

	router := netgrpc.NewRouter()

	server := &netgrpc.Server{Router: router}
	server.Listen(netgrpc.SetConf("the/path"), handler, middleware)

	// Should add a new route
	assert.NotNil(t, router.PathsWrapper["the/path"], "Should not be nil")

	handlerReply, err := server.HandlerSync(nil, handlerRequest)

	// Should call the handler function
	assert.True(t, callHandler, "Should be true")

	// Should return nil error object
	assert.Nil(t, err, "Should be nil")

	// Should return the right values

	optionsUn := network.NewOptions().Unmarshal(handlerReply.Options)

	assert.Equal(t, "sample reply", optionsUn.GetHeader("Custom-Sample-Reply"),
		"Should be equals to sample reply")
	assert.Equal(t, "sample reply", optionsUn.GetHeader("GR-Sample-Reply"),
		"Should be equals to sample reply")
	assert.Empty(t, optionsUn.GetHeader("HP-Sample-Reply"),
		"Should be empty")
	assert.Empty(t, optionsUn.GetHeader("MQ-Sample-Reply"),
		"Should be empty")
	assert.Equal(t, "content-reply", string(handlerReply.Content),
		"Should be equal to content-reply")
}

func TestServer_Handler(t *testing.T) {
	// Should return an error
	server := &netgrpc.Server{}
	err := server.Handler(nil)
	assert.NotNil(t, err, "Should returns an error")
}

func TestServer_Serve(t *testing.T) {

	// Should not panic
	netgrpc.Listen = func(network, address string) (net.Listener, error) {
		return nil, nil
	}

	netgrpc.NewGRPCServer = func(opt ...grpc.ServerOption) *grpc.Server {
		return nil
	}

	netgrpc.RegisterServeServer = func(s *grpc.Server, srv pb.ServeServer) {
	}

	netgrpc.ReflectionRegister = func(s *grpc.Server) {

	}

	server := netgrpc.NewServer(nil)
	assert.NotPanics(t, func() { server.Serve("address") },
		"Should not panics")

	// Should panic on Listen function failure
	netgrpc.Listen = func(network, address string) (net.Listener, error) {
		return nil, errors.New("error")
	}

	netgrpc.NewGRPCServer = func(opt ...grpc.ServerOption) *grpc.Server {
		return nil
	}

	netgrpc.RegisterServeServer = func(s *grpc.Server, srv pb.ServeServer) {
	}

	netgrpc.ReflectionRegister = func(s *grpc.Server) {
	}

	server = netgrpc.NewServer(nil)
	server.Serve("address")
	errorTracking := <-network.ErrorTrack
	assert.NotNil(t, errorTracking, "Should not be nil")
	fmt.Print(errorTracking)

	// Should panic on Serve function failure

	netgrpc.Listen = func(network, address string) (net.Listener, error) {
		return nil, nil
	}

	netgrpc.NewGRPCServer = func(opt ...grpc.ServerOption) *grpc.Server {
		return nil
	}

	netgrpc.RegisterServeServer = func(s *grpc.Server, srv pb.ServeServer) {
	}

	netgrpc.ReflectionRegister = func(s *grpc.Server) {
	}

	server = netgrpc.NewServer(nil)
	server.Serve("address")
	errorTracking = <-network.ErrorTrack
	assert.NotNil(t, errorTracking, "Should not returns nil")

}
