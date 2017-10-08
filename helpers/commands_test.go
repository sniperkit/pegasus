package helpers_test

import (
	"github.com/cpapidas/pegasus/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commands", func() {

	Describe("Call GetContainerID", func() {

		Context("on success", func() {

			It("should return the container id for valid file path", func() {
				helpers.GetContainerIDScriptPath = "../tests/fixtures/commands/get_container_id.sh"
				containerID := helpers.GetContainerID()
				Expect(containerID).To(BeEquivalentTo("commandSuccess"))
			})

			It("should return the default error message for invalid command", func() {
				helpers.GetContainerIDScriptPath = "invalid command"
				containerID := helpers.GetContainerID()
				Expect(containerID).To(BeEquivalentTo("Container ID not found"))
			})

			It("should return the default error message for undefined container id", func() {
				helpers.GetContainerIDScriptPath = `echo ""`
				containerID := helpers.GetContainerID()
				Expect(containerID).To(BeEquivalentTo("Container ID not found"))
			})

		})

	})

})
