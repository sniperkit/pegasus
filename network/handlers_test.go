package network_test

import (
	"github.com/cpapidas/pegasus/network"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestHandler(t *testing.T) {
	// Should contain a channel param and be type of network.Handler
	handler := func(chanel *network.Channel) {}
	valueOf := reflect.ValueOf(handler).String()
	assert.Equal(t, "<func(*network.Channel) Value>", valueOf, "Should be type of network.Handler")
}

func TestMiddleware(t *testing.T) {
	// Should contain a channel param and be type of network.Middleware
	middleware := func(handler network.Handler, channel *network.Channel) {}
	valueOf := reflect.ValueOf(middleware).String()
	assert.Equal(t, "<func(network.Handler, *network.Channel) Value>", valueOf,
		"Should be type of network.Middleware")
}
