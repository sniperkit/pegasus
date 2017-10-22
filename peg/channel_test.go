package peg_test

import (
	"github.com/cpapidas/pegasus/peg"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewChannel(t *testing.T) {
	channel := peg.NewChannel(2000)

	// Should create a new channel and be type of *peg.Channel
	valueOf := reflect.ValueOf(channel).String()
	assert.Equal(t, "<*peg.Channel Value>", valueOf, "Should be type of *peg.Channel")

	// Should send and receive though channel
	channel.Send(peg.BuildPayload([]byte("data"), []byte("options")))
	data := channel.Receive()
	assert.Equal(t, []byte("data"), data.Body, "Should get the Body field and be equal to data")
	assert.Equal(t, []byte("options"), data.Options, "Should get the Options and be equal to options")
}
