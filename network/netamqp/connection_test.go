package netamqp_test

import (
	"bitbucket.org/code_horse/pegasus/network/netamqp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/streadway/amqp"
)

var _ = Describe("Connection", func() {

	Describe("Test Connection struct", func() {

		Context("Test constructor", func() {

			var called bool

			BeforeEach(func() {
				// Set the mocked variable back to originals
				netamqp.Dial = amqp.Dial
				netamqp.NewConnection = NewConnection

				netamqp.RetriesTimes = 1
				netamqp.Sleep = 1

				called = false

				netamqp.Dial = func(url string) (*amqp.Connection, error) {
					called = true
					return &amqp.Connection{}, nil
				}
			})

			It("Should return a object", func() {
				connection, err := netamqp.NewConnection("")
				Expect(connection).ToNot(BeNil())
				Expect(err).To(BeNil())
				Expect(called).To(BeTrue())
			})

		})

	})

})
