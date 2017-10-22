package netgrpc_test

import (
	"github.com/cpapidas/pegasus/network/netgrpc"

	"context"
	"errors"
	"github.com/cpapidas/pegasus/network"
	pb "github.com/cpapidas/pegasus/network/netgrpc/proto"
	"github.com/cpapidas/pegasus/tests/mocks/mnetgrpc"
	"google.golang.org/grpc"
	"testing"
	"github.com/stretchr/testify/assert"

	"reflect"
)

var client = netgrpc.NewClient

func setClientDefaults() {
	netgrpc.NewServerClient = pb.NewServeClient
	netgrpc.Dial = grpc.Dial
	netgrpc.NewClient = client
	netgrpc.RetriesTimes = 1
	netgrpc.Sleep = 0
}

func TestNewClient(t *testing.T) {
	setClientDefaults()

	 //Should return a nil client object
	netgrpc.Dial = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}
	client := netgrpc.NewClient("")
	assert.NotNil(t, client, "Should not be nil")

	// Should be type of *netgrpc.Client
	netgrpc.Dial = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}
	client = netgrpc.NewClient("")
	assert.Equal(t, "<*netgrpc.Client Value>", reflect.ValueOf(client).String(),
		"Should be type of <*netgrpc.Client Value>")

	// Should return an error on Dial failure
	netgrpc.Dial = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return nil, errors.New("error")
	}
	netgrpc.NewServerClient = func(cc *grpc.ClientConn) pb.ServeClient {
		return nil
	}
	client = netgrpc.NewClient("")
	_, err := client.Send(netgrpc.SetConf("path"), network.BuildPayload(nil, nil))
	assert.NotNil(t, err, "Should return an error")
}

func TestClient_Send(t *testing.T) {
	setClientDefaults()

	// Mock the connection
	clientConnection := &mnetgrpc.MockServeClient{}

	// Mock the handler
	clientConnection.HandlerSyncMock = func(
		ctx context.Context,
		in *pb.HandlerRequest,
		opts ...grpc.CallOption,
	) (*pb.HandlerReply, error) {
		return &pb.HandlerReply{
			Content: []byte("content"),
			Options: []byte("options"),
		}, nil
	}

	// Replace the NewServerClient in order to return the mocked object
	netgrpc.NewServerClient = func(cc *grpc.ClientConn) pb.ServeClient {
		return clientConnection
	}

	// Replace the Dial in order to return the connection
	netgrpc.Dial = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}

	client := netgrpc.NewClient("/wherever/whenever/We're/meant/to/be/together/Shakira")

	payload, err := client.Send(netgrpc.SetConf("/Lucky/you/were/born"), network.Payload{})

	// Should return a nil error
	assert.Nil(t,  err, "Should not return an error")


	// Should return a payload.Content equals to content
	assert.Equal(t, "content", string(payload.Body), "Should be equals to content")


	// Should return a payload.Options equals to options
	assert.Equal(t, "options", string(payload.Options), "Should be equals to options")
}

func TestClient_Send_failure(t *testing.T) {
	setClientDefaults()

	// Mock the connection
	clientConnection := &mnetgrpc.MockServeClient{}

	// Mock the handler
	clientConnection.HandlerSyncMock = func(
		ctx context.Context,
		in *pb.HandlerRequest,
		opts ...grpc.CallOption,
	) (*pb.HandlerReply, error) {
		return nil, errors.New("error")
	}

	// Replace the NewServerClient in order to return the mocked object
	netgrpc.NewServerClient = func(cc *grpc.ClientConn) pb.ServeClient {
		return clientConnection
	}

	// Replace the Dial in order to return the connection
	netgrpc.Dial = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}

	client := netgrpc.NewClient("/wherever/whenever/We're/meant/to/be/together/Shakira")

	_, err := client.Send(netgrpc.SetConf("/Lucky/you/were/born"), network.Payload{})

	// Should not return a nil error
	assert.NotNil(t,  err, "Should return an error")
}

func TestClient_Send_nilConnection(t *testing.T) {
	setClientDefaults()

	// Replace the NewServerClient in order to return the mocked object
	netgrpc.NewServerClient = func(cc *grpc.ClientConn) pb.ServeClient {
		return nil
	}

	// Replace the Dial in order to return the connection
	netgrpc.Dial = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}

	client := netgrpc.NewClient("/wherever/whenever/We're/meant/to/be/together/Shakira")

	_, err := client.Send(netgrpc.SetConf("/Lucky/you/were/born"), network.Payload{})

	// Should not return a nil error
	assert.NotNil(t, err, "Should return an error")
}

func TestClient_Close(t *testing.T) {
	setClientDefaults()

	// Mock the client connection
	mockClientConnection := &mnetgrpc.MockClientConnection{}

	// Mock the Close function
	mockClientConnection.CloseMock = func() error {
		return errors.New("string")
	}

	client := &netgrpc.Client {
		Connection: mockClientConnection,
	}

	err := client.Close()

	// Should not return a nil error
	assert.NotNil(t,  err, "Should return an error")

	// Mock the client connection
	mockClientConnection = &mnetgrpc.MockClientConnection{}

	// Mock the Close function
	mockClientConnection.CloseMock = func() error {
		return nil
	}

	client = &netgrpc.Client {
		Connection: mockClientConnection,
	}

	err = client.Close()

	// Should not return an error
	assert.Nil(t,  err, "Should not return an error")
}