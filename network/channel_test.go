package network_test

import (
	"bitbucket.org/code_horse/pegasus/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Channel", func() {

	Describe("Channel struct", func() {

		Context("Channel Create new Chanel function test", func() {

			It("should create a new channel", func() {
				channel := network.NewChannel(2000)
				Expect(reflect.ValueOf(channel).String()).To(Equal("<*network.Channel Value>"))
			})

		})

		Context("Channel Send & Receive function test", func() {

			It("Should send the data to the channel", func() {
				channel := network.NewChannel(2000)
				channel.Send(network.BuildPayload([]byte("data"), []byte("options")))
				data := channel.Receive()
				Expect(data.Body).To(Equal([]byte("data")))
				Expect(data.Options).To(Equal([]byte("options")))
			})

		})

	})

})
