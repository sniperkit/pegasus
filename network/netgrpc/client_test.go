package netgrpc_test

import (
	"bitbucket.org/code_horse/pegasus/network/netgrpc"

	"bitbucket.org/code_horse/pegasus/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Client", func() {

	Describe("Client struct", func() {

		Context("Constructor", func() {
			client := netgrpc.NewClient("")

			It("Should not be nil", func() {
				Expect(client).ToNot(BeNil())
			})

			It("Should be type of *netgrpc.Client", func() {
				Expect(reflect.ValueOf(client).String()).To(Equal("<*netgrpc.Client Value>"))
			})
		})

		Context("Send function", func() {
			client := netgrpc.NewClient("")

			It("Should not panic", func() {
				Expect(func() { client.Send([]string{"path"}, network.BuildPayload(nil, nil)) }).ToNot(Panic())
			})

			It("Should return a type of payload and a error", func() {
				payload, err := client.Send([]string{"path"}, network.BuildPayload(nil, nil))
				Expect(reflect.ValueOf(payload).String()).To(Equal("<*network.Payload Value>"))
				Expect(reflect.ValueOf(err).String()).To(Equal("<*status.statusError Value>"))
			})
		})

	})

})
