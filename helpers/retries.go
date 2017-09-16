package helpers

import (
	"time"
)

// Retries function execute the callable function that takes as parameter for X times with sleep.
// The important thing that you have to know is that callable parameter must to be boolean and
// returns true if wants to continue running. An example could be Retries(4, 3, func() bool { return true })
// the function above will run 4 times but Retries(4, 3, func() bool { return false }) will run only once.
func Retries(times int, seconds int, callback func(...interface{}) bool, params ...interface{}) {
	for i := 0; i < times; i++ {
		if callback(params...) {
			time.Sleep(time.Second * time.Duration(seconds))
		} else {
			return
		}
	}
}
