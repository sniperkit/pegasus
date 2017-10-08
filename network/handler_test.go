package network_test

import (
	"github.com/cpapidas/pegasus/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Handler", func() {

	Describe("Handler struct", func() {

		Context("Handler struct properties", func() {

			It("Should have a channel property", func() {

				var handler network.Handler = func(chanel *network.Channel) {}
				Expect(reflect.ValueOf(handler).String()).To(Equal("<network.Handler Value>"))

			})

		})

	})

})
