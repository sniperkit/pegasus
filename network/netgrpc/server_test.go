package netgrpc_test

import (
	"bitbucket.org/code_horse/pegasus/network/netgrpc"

	"bitbucket.org/code_horse/pegasus/tests/mocks/mockgrpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
	"bitbucket.org/code_horse/pegasus/network"
)

var _ = Describe("Server", func() {

	Describe("Should configure the server properly", func() {

		Context("NewServer function", func() {

			server := netgrpc.NewServer(nil)

			It("Should return a new server instance which implement the IServer interface", func() {
				Expect(server).ToNot(BeNil())
				Expect(reflect.ValueOf(server).String()).To(Equal("<*netgrpc.Server Value>"))
			})

		})

		Context("Serve struct function", func() {

			netgrpc.NewServer = mockgrpc.NewServer

			server := netgrpc.NewServer(nil)

			It("Should listen to Serve function", func() {
				Expect(func() { server.Serve("x.x.x.x:PP") }).ToNot(Panic())
			})

		})

		Context("Listen struct function", func() {

			netgrpc.NewServer = mockgrpc.NewServer

			server := netgrpc.NewServer(nil)

			It("Should have the listen function", func() {
				Expect(func() { server.Listen("/grpc", handler, middleware) }).ToNot(Panic())
			})

		})

		Context("ListenTo struct function", func() {

			netgrpc.NewServer = mockgrpc.NewServer

			server := netgrpc.NewServer(nil)

			It("Should listen to Serve function", func() {
				Expect(func() { server.ListenTo(&netgrpc.Properties{}, handler, middleware) }).ToNot(Panic())
			})

		})

	})

})

func handler(channel *network.Channel) {

}

func middleware(handler network.Handler, chanel *network.Channel){

}