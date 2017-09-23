package netgrpc_test

import (
	"bitbucket.org/code_horse/pegasus/network/netgrpc"

	"bitbucket.org/code_horse/pegasus/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("PathWrapper", func() {

	Describe("PathWrapper struct", func() {

		It("Should returns a new object with valid properties", func() {
			pathWrapper := netgrpc.PathWrapper{}
			pathWrapper.Handler = func(chanel *network.Channel) {}
			pathWrapper.Middleware = func(handler network.Handler, chanel *network.Channel) {}
			Expect(pathWrapper).ToNot(BeNil())
			Expect(reflect.ValueOf(pathWrapper).String()).To(Equal("<netgrpc.PathWrapper Value>"))
			Expect(reflect.ValueOf(pathWrapper.Handler).String()).To(Equal("<network.Handler Value>"))
			Expect(reflect.ValueOf(pathWrapper.Middleware).String()).To(Equal("<network.Middleware Value>"))
		})

	})

})
