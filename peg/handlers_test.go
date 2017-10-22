package peg_test

import (
	"github.com/cpapidas/pegasus/peg"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestHandler(t *testing.T) {
	// Should contain a channel param and be type of peg.Handler
	handler := func(chanel *peg.Channel) {}
	valueOf := reflect.ValueOf(handler).String()
	assert.Equal(t, "<func(*peg.Channel) Value>", valueOf, "Should be type of peg.Handler")
}

func TestMiddleware(t *testing.T) {
	// Should contain a channel param and be type of peg.Middleware
	middleware := func(handler peg.Handler, channel *peg.Channel) {}
	valueOf := reflect.ValueOf(middleware).String()
	assert.Equal(t, "<func(peg.Handler, *peg.Channel) Value>", valueOf,
		"Should be type of peg.Middleware")
}
