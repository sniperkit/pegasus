package helpers_test

import (
	"github.com/cpapidas/pegasus/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsHTTPValidHeader(t *testing.T) {
	// Should return true for valid header
	assert.True(t, helpers.IsHTTPValidHeader("whatever"), "Should return true for valid header")
	//Should return false for invalid header
	assert.False(t, helpers.IsHTTPValidHeader("MQ-FF"), "Should return false for invalid header")
	assert.False(t, helpers.IsHTTPValidHeader("GR-FF"), "Should return false for invalid header")
}

func TestIsGRPCValidHeader(t *testing.T) {
	// Should return a true for valid header
	assert.True(t, helpers.IsGRPCValidHeader("whatever"), "Should return true for valid header")
	// Should return false for invalid header
	assert.False(t, helpers.IsGRPCValidHeader("HP-FF"), "Should return false for invalid header")
	assert.False(t, helpers.IsGRPCValidHeader("MQ-FF"), "Should return false for invalid header")
}

func TestIsAMQPValidHeader(t *testing.T) {
	// Should return a true for valid header
	assert.True(t, helpers.IsAMQPValidHeader("whatever"), "Should return true for valid header")
	// Should return false for invalid header
	assert.False(t, helpers.IsAMQPValidHeader("HP-FF"), "Should return false for invalid header")
	assert.False(t, helpers.IsAMQPValidHeader("GR-FF"), "Should return false for invalid header")
}

func TestAMQPParam(t *testing.T) {
	// Should return the param key name for a valid param
	assert.Equal(t, helpers.AMQPParam("MP-param"), "param", "Should return the param key name ")
	// Should return the an empty string for invalid param
	assert.Empty(t, helpers.AMQPParam("ffparam"), "Should return the an empty string invalid param")
}
