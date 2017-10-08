package netamqp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cpapidas/pegasus/network/netamqp"

	"testing"
)

var NewConnection = netamqp.NewConnection

func TestNetamqp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Netamqp Suite")
}

