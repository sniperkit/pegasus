package peg

import (
	"time"
)

// Retries function calls the callback function.
//
// The parameters are the "times" which defines how many times the callback function will be called, the "seconds" which
// defines the delay (sleep) between calls, the callback function and the parameters of callback function.
//
// An example could be:
//  Retries(4, 3, func(...interface{}) bool { return true })
//
// the function above will run 4 times but:
//  Retries(4, 3, func(...interface{}) bool { return false })
//
// will run only once.
func Retries(times int, seconds int, callback func(...interface{}) bool, params ...interface{}) {
	for i := 0; i < times; i++ {
		if callback(params...) {
			time.Sleep(time.Second * time.Duration(seconds))
		} else {
			return
		}
	}
}
