package network_test

import (
	"bitbucket.org/code_horse/pegasus/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Options", func() {

	Describe("Option struct", func() {

		Context("Set up struct", func() {

			It("Should returns a new object with valid properties", func() {
				options := network.NewOptions()
				options.Path = "foo"
				options.SetField("foo", "bar", "baz")

				Expect(options).To(PointTo(MatchAllFields(
					Fields{
						"Path":   Equal("foo"),
						"Fields": Equal(map[string]map[string]string{"foo": {"bar": "baz"}}),
					},
				)))
			})

		})

		Context("Check the property Fields, CreateNewField method", func() {

			It("Should panic if field mapper is not stetted", func() {
				options := network.NewOptions()
				options.Path = "foo"
				Expect(func() { options.Fields["foo"]["bar"] = "4" }).To(Panic())
			})

		})

		Context("Check Marshal and Unmarshal methods", func() {

			It("Should Marshal/unmarshal the struct properly", func() {
				options := network.NewOptions()
				options.Path = "a/cool/path"
				options.SetField("foo", "bar", "baz")

				marshaledData := options.Marshal()

				unashamedData := network.NewOptions().Unmarshal(marshaledData)

				Expect(unashamedData).To(PointTo(MatchAllFields(
					Fields{
						"Path":   Equal(options.Path),
						"Fields": Equal(options.Fields),
					},
				)))
			})

		})

		Context("Check Marshal and Unmarshal methods with nil params", func() {

			It("Should unmarshal the struct properly", func() {
				network.NewOptions().Unmarshal(nil)
				Expect(func() { network.NewOptions().Unmarshal(nil) }).ToNot(Panic())
			})

			It("Should marshal the struct properly", func() {
				Expect(func() { network.NewOptions().Marshal() }).ToNot(Panic())
			})

		})

		Context("Set/Get Field", func() {

			It("Should returns always the field value", func() {
				Expect(network.NewOptions().GetField("foo", "faa")).To(Equal(""))
			})

			It("Should set always the field value", func() {
				options := network.NewOptions()
				options.SetField("foo", "faa", "Ga")
				Expect(options.GetField("foo", "faa")).To(Equal("Ga"))
			})

		})

	})

})
