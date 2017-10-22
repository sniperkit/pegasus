package network_test

import (
	"github.com/cpapidas/pegasus/network"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPayload(t *testing.T) {
	// Should returns a payload
	payload := network.NewPayload([]byte("body"), []byte("options"))
	assert.Equal(t, []byte("body"), payload.Body, "Should returns the body")
	assert.Equal(t, []byte("options"), payload.Options, "Should returns the options")
}

func TestBuildPayload(t *testing.T) {
	// Should returns a payload
	payload := network.BuildPayload([]byte("body"), []byte("options"))
	assert.Equal(t, []byte("body"), payload.Body, "Should returns the body")
	assert.Equal(t, []byte("options"), payload.Options, "Should returns the options")
}
