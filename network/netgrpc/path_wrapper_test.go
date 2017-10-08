package netgrpc_test

import (
	"github.com/cpapidas/pegasus/network/netgrpc"

	"github.com/cpapidas/pegasus/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("PathWrapper", func() {

	Describe("PathWrapper struct", func() {

		Context("Struct fields", func() {

			pathWrapper := netgrpc.PathWrapper{}
			pathWrapper.Handler = func(chanel *network.Channel) {}
			pathWrapper.Middleware = func(handler network.Handler, chanel *network.Channel) {}

			It("Should not be nil", func() {
				Expect(pathWrapper).ToNot(BeNil())
			})

			It("Should be type of netgrpc.PathWrapper", func() {
				Expect(reflect.ValueOf(pathWrapper).String()).To(Equal("<netgrpc.PathWrapper Value>"))
			})

			It("Should be type of network.Handler", func() {
				Expect(reflect.ValueOf(pathWrapper.Handler).String()).To(Equal("<network.Handler Value>"))
			})

			It("Should be type of network.Middleware", func() {
				Expect(reflect.ValueOf(pathWrapper.Middleware).String()).
					To(Equal("<network.Middleware Value>"))
			})

		})

	})

})
