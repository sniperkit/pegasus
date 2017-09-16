package network_test

import (
	"bitbucket.org/code_horse/pegasus/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Properties", func() {

	Describe("Properties struct", func() {

		Context("Check the struct", func() {

			It("Should return a new valid object on NewProperties", func() {
				properties := network.NewProperties("/path", "GET")
				properties.CreateNewField("foo")
				properties.Fields["foo"]["bar"] = "baz"
				properties.CreateNewObject("faa")
				properties.Objects["faa"]["bor"] = "boz"
				Expect(properties).To(PointTo(MatchAllFields(
					Fields{
						"Path":    Equal("/path"),
						"Method":  Equal("GET"),
						"Fields":  Equal(map[string]map[string]string{"foo": {"bar": "baz"}}),
						"Objects": Equal(map[string]map[string]interface{}{"faa": {"bor": "boz"}}),
					},
				)))
			})

		})

		Context("Check the property Field, CreateNewField function", func() {

			It("Should throw an error if key is not set", func() {
				properties := network.NewProperties("", "")
				properties.Path = "foo"
				Expect(func() { properties.Fields["foo"]["bar"] = "4" }).To(Panic())
			})

			It("Should set the key normally on CreateNewField", func() {
				properties := network.NewProperties("", "")
				properties.Path = "foo"
				properties.CreateNewField("foo")
				Expect(func() { properties.Fields["foo"]["bar"] = "4" }).ToNot(Panic())
			})

		})

		Context(" Check the property Objects, CreateNewObject function", func() {

			It("Should throw an error if key is not set", func() {
				properties := network.NewProperties("", "")
				properties.Path = "foo"
				Expect(func() { properties.Objects["foo"]["bar"] = "4" }).To(Panic())
			})

			It("Should set the key normally on CreateNewObject", func() {
				properties := network.NewProperties("", "")
				properties.Path = "foo"
				properties.CreateNewObject("foo")
				Expect(func() { properties.Objects["foo"]["bar"] = "4" }).ToNot(Panic())
			})

		})

	})

})
