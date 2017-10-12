package netgrpc_test

import (
	"github.com/cpapidas/pegasus/network/netgrpc"

	"context"
	"errors"
	"github.com/cpapidas/pegasus/network"
	pb "github.com/cpapidas/pegasus/network/netgrpc/proto"
	"github.com/cpapidas/pegasus/tests/mocks/mock_netgrpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"reflect"
)

var _ = Describe("Client", func() {

	Describe("Client struct", func() {

		BeforeEach(func() {
			netgrpc.Dial = grpc.Dial
			netgrpc.RetriesTimes = 1
			netgrpc.Sleep = 0
		})

		Context("Constructor", func() {

			netgrpc.Dial = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
				return &grpc.ClientConn{}, nil
			}

			client := netgrpc.NewClient("")

			It("Should not be nil", func() {
				Expect(client).ToNot(BeNil())
			})

			It("Should be type of *netgrpc.Client", func() {
				Expect(reflect.ValueOf(client).String()).To(Equal("<*netgrpc.Client Value>"))
			})
		})

		Context("Send function", func() {

			clientConnection := &mock_netgrpc.MockClientConnection{}

			clientConnection.HandlerSyncMock = func(
				ctx context.Context,
				in *pb.HandlerRequest,
				opts ...grpc.CallOption,
			) (*pb.HandlerReply, error) {
				return &pb.HandlerReply{
					Content: []byte("content"),
					Options: []byte("options"),
				}, nil
			}

			netgrpc.NewServerClient = func(cc *grpc.ClientConn) pb.ServeClient {
				return clientConnection
			}

			netgrpc.Dial = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
				return &grpc.ClientConn{}, nil
			}

			client := netgrpc.NewClient("/wherever/whenever/We're/meant/to/be/together/Shakira")

			payload, err := client.Send(netgrpc.SetConf("/Lucky/you/were/born"), network.Payload{})

			It("Should return a nil error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return a payload.Content equals to content", func() {
				Expect(string(payload.Body)).To(Equal("content"))
			})

			It("Should return a payload.Options equals to options", func() {
				Expect(string(payload.Options)).To(Equal("options"))
			})
		})

		Context("Send function on failure", func() {

			clientConnection := &mock_netgrpc.MockClientConnection{}

			clientConnection.HandlerSyncMock = func(
				ctx context.Context,
				in *pb.HandlerRequest,
				opts ...grpc.CallOption,
			) (*pb.HandlerReply, error) {
				return nil, errors.New("error")
			}

			netgrpc.NewServerClient = func(cc *grpc.ClientConn) pb.ServeClient {
				return clientConnection
			}

			client := netgrpc.NewClient("/wherever/whenever/We're/meant/to/be/together/Shakira")

			_, err := client.Send(netgrpc.SetConf("/Lucky/you/were/born"), network.Payload{})

			It("Should return a nil error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

	})

})
