package peg_test

import (
	"github.com/cpapidas/pegasus/peg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsHTTPValidHeader(t *testing.T) {
	// Should return true for valid header
	assert.True(t, peg.IsHTTPValidHeader("whatever"), "Should return true for valid header")
	//Should return false for invalid header
	assert.False(t, peg.IsHTTPValidHeader("MQ-FF"), "Should return false for invalid header")
	assert.False(t, peg.IsHTTPValidHeader("GR-FF"), "Should return false for invalid header")
}

func TestIsGRPCValidHeader(t *testing.T) {
	// Should return a true for valid header
	assert.True(t, peg.IsGRPCValidHeader("whatever"), "Should return true for valid header")
	// Should return false for invalid header
	assert.False(t, peg.IsGRPCValidHeader("HP-FF"), "Should return false for invalid header")
	assert.False(t, peg.IsGRPCValidHeader("MQ-FF"), "Should return false for invalid header")
}

func TestIsAMQPValidHeader(t *testing.T) {
	// Should return a true for valid header
	assert.True(t, peg.IsAMQPValidHeader("whatever"), "Should return true for valid header")
	// Should return false for invalid header
	assert.False(t, peg.IsAMQPValidHeader("HP-FF"), "Should return false for invalid header")
	assert.False(t, peg.IsAMQPValidHeader("GR-FF"), "Should return false for invalid header")
}

func TestAMQPParam(t *testing.T) {
	// Should return the param key name for a valid param
	assert.Equal(t, peg.AMQPParam("MP-param"), "param", "Should return the param key name ")
	// Should return the an empty string for invalid param
	assert.Empty(t, peg.AMQPParam("ffparam"), "Should return the an empty string invalid param")
}
