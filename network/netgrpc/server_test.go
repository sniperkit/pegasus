package netgrpc_test

import (
	"bitbucket.org/code_horse/pegasus/network/netgrpc"

	"bitbucket.org/code_horse/pegasus/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Server", func() {

	Describe("Should configure the server properly", func() {

		Context("NewServer function", func() {

			server := netgrpc.NewServer(nil)

			It("Should return a new server instance which implement the Server interface", func() {
				Expect(server).ToNot(BeNil())
				Expect(reflect.ValueOf(server).String()).To(Equal("<*netgrpc.Server Value>"))
			})

		})

		Context("SetPath function", func() {
			path := netgrpc.SetPath("path")

			It("Should return the path as array", func() {
				Expect(path).To(Equal([]string{"path"}))
			})
		})

	})

})

func handler(channel *network.Channel) {

}

func middleware(handler network.Handler, chanel *network.Channel) {

}
