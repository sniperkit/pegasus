package helpers_test

import (
	"github.com/cpapidas/pegasus/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Retries", func() {

	Describe("Call Retries function", func() {

		Context("Test execution and times", func() {

			It("Should run for X times normally", func() {

				var executionTimes = 4
				var times int = 0

				helpers.Retries(executionTimes, 0, func(...interface{}) bool {
					times++
					return true
				}, nil)

				Expect(times).To(Equal(executionTimes))
			})

			It("Should quit if returns false", func() {

				var executionTimes = 4
				var times int = 0

				helpers.Retries(executionTimes, 0, func(...interface{}) bool {
					times++
					if times == 2 {
						return false
					}
					return true
				}, nil)

				Expect(times).To(Equal(2))
			})

			It("Should not called for 0 times", func() {

				var executionTimes = 0
				var times int = 0

				helpers.Retries(executionTimes, 0, func(...interface{}) bool {
					times++
					return true
				}, nil)

				Expect(times).To(Equal(executionTimes))

			})

		})

		Context("Test execution with parameters", func() {

			It("Should get all the params", func() {

				var (
					executionTimes int = 1
					times          int = 0
					param1         string
					param2         string
				)

				helpers.Retries(executionTimes, 0, func(params ...interface{}) bool {
					times++
					param1 = params[0].(string)
					param2 = params[1].(string)
					return true
				}, "params1", "params2")

				Expect(times).To(Equal(executionTimes))
				Expect(param1).To(Equal("params1"))
				Expect(param2).To(Equal("params2"))

			})

		})

	})

})


