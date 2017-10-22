package helpers_test

import (
	"github.com/cpapidas/pegasus/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetries(t *testing.T) {
	var (
		executionTimes = 4
		times          int
	)

	// Should run normally for 4 times
	helpers.Retries(executionTimes, 0, func(...interface{}) bool {
		times++
		return true
	}, nil)
	assert.Equal(t, executionTimes, times, "Should run 4 times")

	// Should quit if returns false and run 2 times
	times = 0
	helpers.Retries(executionTimes, 0, func(...interface{}) bool {
		times++
		if times == 2 {
			return false
		}
		return true
	}, nil)
	assert.Equal(t, 2, times, "Should run 2 times only")

	// Should not called at all
	executionTimes = 0
	times = 0

	helpers.Retries(executionTimes, 0, func(...interface{}) bool {
		times++
		return true
	}, nil)
	assert.Equal(t, executionTimes, times, "Should not called at all")

	// Should get all the params that we pass
	var (
		param1 string
		param2 string
	)
	executionTimes = 1
	times = 0
	helpers.Retries(executionTimes, 0, func(params ...interface{}) bool {
		times++
		param1 = params[0].(string)
		param2 = params[1].(string)
		return true
	}, "params1", "params2")
	assert.Equal(t, executionTimes, times, "Should called 1 time")
	assert.Equal(t, "params1", param1, "Should has the value params1")
	assert.Equal(t, "params2", param2, "Should has the value params2")
}
