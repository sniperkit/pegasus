package helpers_test

import (
	"github.com/cpapidas/pegasus/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Headers", func() {

	Describe("HTTP headers", func() {

		Context("Valida/Invalid Headers", func() {

			It("Should return a true for valid header", func() {
				Expect(helpers.IsHTTPValidHeader("whatever")).To(BeTrue())
			})

			It("Should return false for invalid header", func() {
				Expect(helpers.IsHTTPValidHeader("MQ-FF")).To(BeFalse())
				Expect(helpers.IsHTTPValidHeader("GR-FF")).To(BeFalse())
			})

		})

	})

	Describe("GRPC headers", func() {

		Context("Valida/Invalid Headers", func() {

			It("Should return a true for valid header", func() {
				Expect(helpers.IsGRPCValidHeader("whatever")).To(BeTrue())
			})

			It("Should return false for invalid header", func() {
				Expect(helpers.IsGRPCValidHeader("HP-FF")).To(BeFalse())
				Expect(helpers.IsGRPCValidHeader("MQ-FF")).To(BeFalse())
			})

		})

	})

	Describe("AMQP headers", func() {

		Context("Valida/Invalid Headers", func() {

			It("Should return a true for valid header", func() {
				Expect(helpers.IsAMQPValidHeader("whatever")).To(BeTrue())
			})

			It("Should return false for invalid header", func() {
				Expect(helpers.IsAMQPValidHeader("HP-FF")).To(BeFalse())
				Expect(helpers.IsAMQPValidHeader("GR-FF")).To(BeFalse())
			})

		})

	})

})
