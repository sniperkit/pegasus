package nethttp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNethttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nethttp Suite")
}
