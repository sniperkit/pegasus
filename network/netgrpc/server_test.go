package netgrpc_test

import (
	"bitbucket.org/code_horse/pegasus/network/netgrpc"

	"bitbucket.org/code_horse/pegasus/network"
	pb "bitbucket.org/code_horse/pegasus/network/netgrpc/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

		//Context("HandleSync function with middleware", func() {
		//
		//	options := network.NewOptions()
		//
		//	options.SetHeader("Custom-Sample", "sample")
		//	options.SetHeader("GR-Sample", "sample")
		//	options.SetHeader("HP-Sample", "sample")
		//	options.SetHeader("MQ-Sample", "sample")
		//
		//	handlerRequest := &pb.HandlerRequest{
		//		Content: []byte("content"),
		//		Options: options.Marshal(),
		//		Path:    "the/path",
		//	}
		//
		//	var handler = func(channel *network.Channel) {}
		//	var middleware = func(handler network.Handler, channel *network.Channel) {}
		//
		//	router := netgrpc.NewRouter()
		//
		//	server := netgrpc.NewServer(router)
		//	server.Listen(netgrpc.SetConf("whatever"), handler, middleware)
		//
		//	It("Should add a new route", func() {
		//		Expect(router.PathsWrapper["whatever"]).ToNot(BeNil())
		//	})
		//
		//})

	})

})
