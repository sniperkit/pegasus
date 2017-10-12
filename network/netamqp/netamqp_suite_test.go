package netamqp_test

import (
	"github.com/cpapidas/pegasus/network/netamqp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var NewConnection = netamqp.NewConnection

func TestNetamqp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Netamqp Suite")
}
