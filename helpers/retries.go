package helpers

import (
	"time"
)

func Retries(times int, seconds int, callback func(...interface{}) bool, params ...interface{}) {
	for i := 0; i < times; i++ {
		if callback(params) {
			time.Sleep(time.Second * time.Duration(seconds))
		} else {
			return
		}
	}
}
