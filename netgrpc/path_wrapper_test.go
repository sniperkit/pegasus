package netgrpc_test

import (
	"github.com/cpapidas/pegasus/netgrpc"
	"github.com/cpapidas/pegasus/peg"
	"reflect"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPathWrapper(t *testing.T) {

	pathWrapper := netgrpc.PathWrapper{}
	pathWrapper.Handler = func(chanel *peg.Channel) {}
	pathWrapper.Middleware = func(handler peg.Handler, chanel *peg.Channel) {}

	// Should not be nil
	assert.NotNil(t, pathWrapper, "Should not be nil")

	// Should be type of netgrpc.PathWrapper
	assert.Equal(t, "<netgrpc.PathWrapper Value>", reflect.ValueOf(pathWrapper).String(),
		"Should be equals to <netgrpc.PathWrapper Value>")

	// Should be type of peg.Handler
	assert.Equal(t, "<peg.Handler Value>", reflect.ValueOf(pathWrapper.Handler).String(),
		"Should be equals to <peg.Handler Value>")

	// Should be type of peg.Middleware
	assert.Equal(t, "<peg.Middleware Value>", reflect.ValueOf(pathWrapper.Middleware).String(),
		"Should be equals to <peg.Middleware Value>")
}
