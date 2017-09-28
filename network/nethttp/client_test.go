package nethttp_test

import (
	"bitbucket.org/code_horse/pegasus/network/nethttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Client", func() {

	Describe("Client struct", func() {

		Context("Constructor", func() {
			client := nethttp.NewClient()

			It("Should not be nil", func() {
				Expect(client).ToNot(BeNil())
			})

			It("Should be type of *nethttp.Client", func() {
				Expect(reflect.ValueOf(client).String()).To(Equal("<*nethttp.Client Value>"))
			})
		})

	})

})
