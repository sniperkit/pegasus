package blunder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"errors"
	"bitbucket.org/code_horse/pegasus/blunder"
)

var _ = Describe("Error", func() {

	Describe("Generate an error", func() {

		Context("When handle method", func() {

			It("Should silence return on error param nil", func(){
				err := blunder.Set("This is an error message", nil)
				err.Handle()
				Expect(err.Handle).ToNot(Panic())
			})

		})

		Context("When error importance is 4 (PanicError)", func() {

			It("Should panic and stop the execution", func() {
				err := blunder.Set("This is the message", errors.New("sample error"))
				err.SetPanicError()
				Expect(err.Handle).To(Panic())
			})

		})

	})

})
