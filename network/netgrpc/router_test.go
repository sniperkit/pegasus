package netgrpc_test

import (
	"github.com/cpapidas/pegasus/network/netgrpc"

	"github.com/cpapidas/pegasus/network"
	"reflect"
	"testing"
	"github.com/stretchr/testify/assert"

)

func TestNewRouter(t *testing.T) {

	router := netgrpc.NewRouter()

	// Should not be nil
	assert.NotNil(t, router, "Should not be nil")

	// Should be type of *netgrpc.Router
	assert.Equal(t, "<*netgrpc.Router Value>", reflect.ValueOf(router).String(),
		"Should be type of <*netgrpc.Router Value>")

}

func TestRouter_Add(t *testing.T) {

	router := netgrpc.NewRouter()

	handler := func(channel *network.Channel) {}

	middleware := func(handler network.Handler, channel *network.Channel) {}

	router.Add("/goo/gaa", handler, middleware)
	// Should add a new PathsWrapper
	assert.NotNil(t, router.PathsWrapper["/goo/gaa"], "Should not be nil")
}
