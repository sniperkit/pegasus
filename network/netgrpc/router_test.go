package netgrpc_test

import (
	"bitbucket.org/code_horse/pegasus/network/netgrpc"

	"bitbucket.org/code_horse/pegasus/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Router", func() {

	Describe("Router struct", func() {

		Context("Router constructor", func() {

			router := netgrpc.NewRouter()

			It("Should not be nil", func() {
				Expect(router).ToNot(BeNil())
			})

			It("Should be type of *netgrpc.Router", func() {
				Expect(reflect.ValueOf(router).String()).To(Equal("<*netgrpc.Router Value>"))
			})

		})

		Context("Router Add function", func() {

			router := netgrpc.NewRouter()

			handler := func(channel *network.Channel) {}

			middleware := func(handler network.Handler, channel *network.Channel) {}

			It("Should add a new PathsWrapper", func() {
				router.Add("/goo/gaa", handler, middleware)
			})

		})

	})

})
