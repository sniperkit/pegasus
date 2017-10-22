package network_test

import (
	"github.com/cpapidas/pegasus/network"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewChannel(t *testing.T) {
	channel := network.NewChannel(2000)

	// Should create a new channel and be type of *network.Channel
	valueOf := reflect.ValueOf(channel).String()
	assert.Equal(t, "<*network.Channel Value>", valueOf, "Should be type of *network.Channel")

	// Should send and receive though channel
	channel.Send(network.BuildPayload([]byte("data"), []byte("options")))
	data := channel.Receive()
	assert.Equal(t, []byte("data"), data.Body, "Should get the Body field and be equal to data")
	assert.Equal(t, []byte("options"), data.Options, "Should get the Options and be equal to options")
}
