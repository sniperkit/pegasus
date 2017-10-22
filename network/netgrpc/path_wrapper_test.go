package netgrpc_test

import (
	"github.com/cpapidas/pegasus/network/netgrpc"
	"github.com/cpapidas/pegasus/network"
	"reflect"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPathWrapper(t *testing.T) {

	pathWrapper := netgrpc.PathWrapper{}
	pathWrapper.Handler = func(chanel *network.Channel) {}
	pathWrapper.Middleware = func(handler network.Handler, chanel *network.Channel) {}

	// Should not be nil
	assert.NotNil(t, pathWrapper, "Should not be nil")

	// Should be type of netgrpc.PathWrapper
	assert.Equal(t, "<netgrpc.PathWrapper Value>", reflect.ValueOf(pathWrapper).String(),
		"Should be equals to <netgrpc.PathWrapper Value>")

	// Should be type of network.Handler
	assert.Equal(t, "<network.Handler Value>", reflect.ValueOf(pathWrapper.Handler).String(),
		"Should be equals to <network.Handler Value>")

	// Should be type of network.Middleware
	assert.Equal(t, "<network.Middleware Value>", reflect.ValueOf(pathWrapper.Middleware).String(),
		"Should be equals to <network.Middleware Value>")
}
