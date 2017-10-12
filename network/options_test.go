package network_test

import (
	"github.com/cpapidas/pegasus/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Options", func() {

	Describe("Option struct", func() {

		Context("Set up struct", func() {

			It("Should returns a new object with valid properties", func() {
				options := network.NewOptions()
				options.SetField("foo", "bar", "baz")

				Expect(options).To(PointTo(MatchAllFields(
					Fields{
						"Fields": Equal(map[string]map[string]string{"foo": {"bar": "baz"}}),
					},
				)))
			})

		})

		Context("Check the property Fields, CreateNewField method", func() {

			It("Should panic if field mapper is not stetted", func() {
				options := network.NewOptions()
				Expect(func() { options.Fields["foo"]["bar"] = "4" }).To(Panic())
			})

		})

		Context("Check Marshal and Unmarshal methods", func() {

			It("Should Marshal/unmarshal the struct properly", func() {
				options := network.NewOptions()
				options.SetField("foo", "bar", "baz")

				marshaledData := options.Marshal()

				unashamedData := network.NewOptions().Unmarshal(marshaledData)

				Expect(unashamedData).To(PointTo(MatchAllFields(
					Fields{
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

		Context("Set/Get Params", func() {

			options := network.NewOptions()

			options.SetParams(map[string]string{"foo": "fa"})

			It("Should be equal to mapper", func() {
				Expect(options.Fields["PARAMS"]).To(Equal(map[string]string{"foo": "fa"}))
			})

			It("Should return map[string]string", func() {
				Expect(options.GetParams()).To(Equal(map[string]string{"foo": "fa"}))
			})

			It("Should set a new param", func() {
				options.SetParam("baz", "ba")
				Expect(options.GetParams()).To(Equal(map[string]string{"foo": "fa", "baz": "ba"}))
			})

			It("Should get a param", func() {
				Expect(options.GetParam("foo")).To(Equal("fa"))
			})
		})

		Context("Set/Get Headers", func() {

			options := network.NewOptions()

			options.SetHeaders(map[string]string{"foo": "fa"})

			It("Should be equal to mapper", func() {
				Expect(options.Fields["HEADERS"]).To(Equal(map[string]string{"foo": "fa"}))
			})

			It("Should return map[string]string", func() {
				Expect(options.GetHeaders()).To(Equal(map[string]string{"foo": "fa"}))
			})

			It("Should set a new header", func() {
				options.SetHeader("baz", "ba")
				Expect(options.GetHeaders()).To(Equal(map[string]string{"foo": "fa", "baz": "ba"}))
			})

			It("Should get a header", func() {
				Expect(options.GetHeader("foo")).To(Equal("fa"))
			})
		})

	})

})
