package netgrpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNetgrpc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Netgrpc Suite")
}
