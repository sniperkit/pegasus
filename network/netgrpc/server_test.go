package netgrpc_test

import (
	"github.com/cpapidas/pegasus/network/netgrpc"

	"errors"
	"fmt"
	"github.com/cpapidas/pegasus/network"
	pb "github.com/cpapidas/pegasus/network/netgrpc/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"net"
	"reflect"
)

var _ = Describe("Server", func() {

	Describe("Should configure the server properly", func() {

		Context("NewServer function", func() {

			server := netgrpc.NewServer(nil)

			It("Should return a new server instance which implement the Server interface", func() {
				Expect(server).ToNot(BeNil())
				Expect(reflect.ValueOf(server).String()).To(Equal("<*netgrpc.Server Value>"))
			})

		})

		Context("SetConf function", func() {
			path := netgrpc.SetConf("path")

			It("Should return the path as array", func() {
				Expect(path).To(Equal([]string{"path"}))
			})
		})

		Context("Listen function", func() {

			var handler = func(channel *network.Channel) {}
			var middleware = func(handler network.Handler, channel *network.Channel) {}

			router := netgrpc.NewRouter()

			server := netgrpc.NewServer(router)
			server.Listen(netgrpc.SetConf("whatever"), handler, middleware)

			It("Should add a new route", func() {
				Expect(router.PathsWrapper["whatever"]).ToNot(BeNil())
			})

		})

		Context("HandleSync function", func() {

			callHandler := false

			options := network.NewOptions()

			options.SetHeader("GR-Sample", "sample")
			options.SetHeader("HP-Sample", "sample")
			options.SetHeader("MQ-Sample", "sample")
			options.SetHeader("Custom-Sample", "sample")

			handlerRequest := &pb.HandlerRequest{
				Content: []byte("content"),
				Options: options.Marshal(),
				Path:    "the/path",
			}

			var handler = func(channel *network.Channel) {

				callHandler = true

				payload := channel.Receive()
				options := network.NewOptions().Unmarshal(payload.Options)

				It("Should receive valid params", func() {
					Expect(options.GetHeader("Custom-Sample")).To(Equal("sample"))
					Expect(options.GetHeader("GR-Sample")).To(Equal("sample"))
					Expect(options.GetHeader("HP-Sample")).To(BeEmpty())
					Expect(options.GetHeader("MQ-Sample")).To(BeEmpty())
					Expect(string(payload.Body)).To(Equal("content"))
				})

				replyOptions := network.NewOptions()
				replyOptions.SetHeader("Custom-Sample-Reply", "sample reply")
				replyOptions.SetHeader("GR-Sample-Reply", "sample reply")
				replyOptions.SetHeader("HP-Sample-Reply", "sample reply")
				replyOptions.SetHeader("MQ-Sample-Reply", "sample reply")

				channel.Send(network.BuildPayload([]byte("content-reply"), replyOptions.Marshal()))
			}

			router := netgrpc.NewRouter()

			server := &netgrpc.Server{Router: router}
			server.Listen(netgrpc.SetConf("the/path"), handler, nil)

			It("Should add a new route", func() {
				Expect(router.PathsWrapper["the/path"]).ToNot(BeNil())
			})

			handlerReply, err := server.HandlerSync(nil, handlerRequest)

			It("Should call the handler function", func() {
				Expect(callHandler).To(BeTrue())
			})

			It("Should return nil error object", func() {
				Expect(err).To(BeNil())
			})

			It("Should return the right values", func() {

				options := network.NewOptions().Unmarshal(handlerReply.Options)

				Expect(options.GetHeader("Custom-Sample-Reply")).To(Equal("sample reply"))
				Expect(options.GetHeader("GR-Sample-Reply")).To(Equal("sample reply"))
				Expect(options.GetHeader("HP-Sample-Reply")).To(BeEmpty())
				Expect(options.GetHeader("MQ-Sample-Reply")).To(BeEmpty())
				Expect(string(handlerReply.Content)).To(Equal("content-reply"))
			})

		})

		Context("HandleSync function midleware", func() {

			callHandler := false

			options := network.NewOptions()

			options.SetHeader("GR-Sample", "sample")
			options.SetHeader("HP-Sample", "sample")
			options.SetHeader("MQ-Sample", "sample")
			options.SetHeader("Custom-Sample", "sample")

			handlerRequest := &pb.HandlerRequest{
				Content: []byte("content"),
				Options: options.Marshal(),
				Path:    "the/path",
			}

			var handler = func(channel *network.Channel) {

				callHandler = true

				payload := channel.Receive()
				options := network.NewOptions().Unmarshal(payload.Options)

				It("Should receive valid params", func() {
					Expect(options.GetHeader("Custom-Sample")).To(Equal("sample reply middleware"))
					Expect(options.GetHeader("GR-Sample")).To(Equal("sample reply middleware"))
					Expect(string(payload.Body)).To(Equal("content"))
				})

				replyOptions := network.NewOptions()
				replyOptions.SetHeader("Custom-Sample-Reply", "sample reply")
				replyOptions.SetHeader("GR-Sample-Reply", "sample reply")
				replyOptions.SetHeader("HP-Sample-Reply", "sample reply")
				replyOptions.SetHeader("MQ-Sample-Reply", "sample reply")

				channel.Send(network.BuildPayload([]byte("content-reply"), replyOptions.Marshal()))
			}

			var middleware = func(handler network.Handler, channel *network.Channel) {

				callHandler = true

				payload := channel.Receive()
				options := network.NewOptions().Unmarshal(payload.Options)

				It("Should receive valid params", func() {
					Expect(options.GetHeader("Custom-Sample")).To(Equal("sample"))
					Expect(options.GetHeader("GR-Sample")).To(Equal("sample"))
					Expect(options.GetHeader("HP-Sample")).To(BeEmpty())
					Expect(options.GetHeader("MQ-Sample")).To(BeEmpty())
					Expect(string(payload.Body)).To(Equal("content"))
				})

				replyOptions := network.NewOptions()
				replyOptions.SetHeader("Custom-Sample", "sample reply middleware")
				replyOptions.SetHeader("GR-Sample", "sample reply middleware")

				channel.Send(network.BuildPayload([]byte("content"), replyOptions.Marshal()))
				handler(channel)
			}

			router := netgrpc.NewRouter()

			server := &netgrpc.Server{Router: router}
			server.Listen(netgrpc.SetConf("the/path"), handler, middleware)

			It("Should add a new route", func() {
				Expect(router.PathsWrapper["the/path"]).ToNot(BeNil())
			})

			handlerReply, err := server.HandlerSync(nil, handlerRequest)

			It("Should call the handler function", func() {
				Expect(callHandler).To(BeTrue())
			})

			It("Should return nil error object", func() {
				Expect(err).To(BeNil())
			})

			It("Should return the right values", func() {

				options := network.NewOptions().Unmarshal(handlerReply.Options)

				Expect(options.GetHeader("Custom-Sample-Reply")).To(Equal("sample reply"))
				Expect(options.GetHeader("GR-Sample-Reply")).To(Equal("sample reply"))
				Expect(options.GetHeader("HP-Sample-Reply")).To(BeEmpty())
				Expect(options.GetHeader("MQ-Sample-Reply")).To(BeEmpty())
				Expect(string(handlerReply.Content)).To(Equal("content-reply"))
			})

		})

		Context("Handler function", func() {

			It("Should return an error", func() {
				server := &netgrpc.Server{}
				err := server.Handler(nil)
				Expect(err).ToNot(BeNil())
			})

		})

		Context("Serve function", func() {

			It("Should not panic", func() {

				netgrpc.Listen = func(network, address string) (net.Listener, error) {
					return nil, nil
				}

				netgrpc.NewGRPCServer = func(opt ...grpc.ServerOption) *grpc.Server {
					return nil
				}

				netgrpc.RegisterServeServer = func(s *grpc.Server, srv pb.ServeServer) {
				}

				netgrpc.ReflectionRegister = func(s *grpc.Server) {

				}

				server := netgrpc.NewServer(nil)
				Expect(func() { server.Serve("address") }).ToNot(Panic())
			})

			It("Should panic on Listen function failure", func(done Done) {

				netgrpc.Listen = func(network, address string) (net.Listener, error) {
					return nil, errors.New("error")
				}

				netgrpc.NewGRPCServer = func(opt ...grpc.ServerOption) *grpc.Server {
					return nil
				}

				netgrpc.RegisterServeServer = func(s *grpc.Server, srv pb.ServeServer) {
				}

				netgrpc.ReflectionRegister = func(s *grpc.Server) {
				}

				server := netgrpc.NewServer(nil)
				server.Serve("address")
				errorTracking := <-network.ErrorTrack
				Expect(errorTracking).ToNot(BeNil())
				fmt.Print(errorTracking)
				close(done)
			})

			It("Should panic on Serve function failure", func(done Done) {

				netgrpc.Listen = func(network, address string) (net.Listener, error) {
					return nil, nil
				}

				netgrpc.NewGRPCServer = func(opt ...grpc.ServerOption) *grpc.Server {
					return nil
				}

				netgrpc.RegisterServeServer = func(s *grpc.Server, srv pb.ServeServer) {
				}

				netgrpc.ReflectionRegister = func(s *grpc.Server) {
				}

				server := netgrpc.NewServer(nil)
				server.Serve("address")
				errorTracking := <-network.ErrorTrack
				Expect(errorTracking).ToNot(BeNil())
				close(done)
			})

		})

	})

})
