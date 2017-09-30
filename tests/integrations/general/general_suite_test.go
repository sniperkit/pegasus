package general_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGeneral(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "General Suite")
}
