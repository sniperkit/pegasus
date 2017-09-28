package nethttp_test

import (
	"bitbucket.org/code_horse/pegasus/network/nethttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Server", func() {

	Describe("Server struct", func() {

		Context("Construct NewServer", func() {

			It("Should not be nil", func() {
				server := nethttp.NewServer(nil)
				Expect(server).ToNot(BeNil())
			})

			It("Should be type of *Server", func() {
				server := nethttp.NewServer(nil)
				Expect(reflect.ValueOf(server).String()).To(Equal("<*nethttp.Server Value>"))
			})

		})

		Context("SetPath function", func() {

			It("Should return an array of given strings", func() {
				Expect(nethttp.SetPath("foo", "bar")).To(Equal([]string{"foo", "bar"}))
			})

		})
	})

})
