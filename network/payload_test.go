package network_test

import (
	"bitbucket.org/code_horse/pegasus/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Payload", func() {

	Describe("Payload struct", func() {

		Context("Payload Constructor", func() {

			It("Should returns a payload", func() {
				payload := network.NewPayload([]byte("body"), []byte("options"))
				Expect(payload.Body).To(Equal([]byte("body")))
				Expect(payload.Options).To(Equal([]byte("options")))
			})

		})

		Context("Build a payload struct", func() {

			It("Should returns a payload", func() {
				payload := network.BuildPayload([]byte("body"), []byte("options"))
				Expect(payload.Body).To(Equal([]byte("body")))
				Expect(payload.Options).To(Equal([]byte("options")))
			})

		})

	})

})
