package nethttp_test

import (
	"github.com/cpapidas/pegasus/network/nethttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Method", func() {

	Describe("Method type", func() {

		Context("Const variables", func() {

			It("Should the right const variables", func() {
				Expect(nethttp.Get.String()).To(Equal(http.MethodGet))
				Expect(nethttp.Head.String()).To(Equal(http.MethodHead))
				Expect(nethttp.Post.String()).To(Equal(http.MethodPost))
				Expect(nethttp.Put.String()).To(Equal(http.MethodPut))
				Expect(nethttp.Patch.String()).To(Equal(http.MethodPatch))
				Expect(nethttp.Delete.String()).To(Equal(http.MethodDelete))
				Expect(nethttp.Connect.String()).To(Equal(http.MethodConnect))
				Expect(nethttp.Options.String()).To(Equal(http.MethodOptions))
				Expect(nethttp.Trace.String()).To(Equal(http.MethodTrace))
			})

		})

	})

})
