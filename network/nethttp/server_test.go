package nethttp_test

import (
	"github.com/cpapidas/pegasus/network/nethttp"

	"bytes"
	"github.com/cpapidas/pegasus/network"
	"github.com/cpapidas/pegasus/tests/mocks/mhttp"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"reflect"
)

var _ = Describe("Server", func() {

	Describe("Server struct", func() {

		Context("Construct NewServer", func() {

			It("Should not be nil", func() {
				server := nethttp.NewServer(nil)
				Expect(server).ToNot(BeNil())
			})

			It("Should be type of *Server", func() {
				server := nethttp.NewServer(nil)
				Expect(reflect.ValueOf(server).String()).To(Equal("<*nethttp.Server Value>"))
			})

		})

		Context("SetConf function", func() {

			It("Should return an array of given strings", func() {
				Expect(nethttp.SetConf("foo", "bar")).To(Equal([]string{"foo", "bar"}))
			})

		})

		Context("Listen function", func() {

			callHandler := false

			router := &mhttp.MockRouter{}

			router.HandleFuncMock = func(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {

				w := &mhttp.MockResponseWriter{
					Headers: make(map[string][]string),
				}

				r, _ := http.NewRequest("POST", "anything", bytes.NewReader([]byte("content")))

				r.Header = make(map[string][]string)
				r.Header["Custom-Sample"] = []string{"Sample"}
				r.Header["HP-Sample"] = []string{"Sample"}
				r.Header["GR-Sample"] = []string{"Sample"}
				r.Header["MQ-Sample"] = []string{"Sample"}

				r.Body = ioutil.NopCloser(bytes.NewReader([]byte("content")))

				f(w, r)

				It("Should contain the right headers", func() {
					Expect(w.Headers["Custom-Sample"]).To(Equal([]string{"sample"}))
					Expect(w.Headers["Hp-Sample"]).To(Equal([]string{"sample"}))
					Expect(w.Headers["Gr-Sample"]).To(BeEmpty())
					Expect(w.Headers["Mq-Sample"]).To(BeEmpty())
					Expect(string(w.Body)).To(Equal("content reply"))
				})

				return &mux.Route{}
			}

			server := nethttp.NewServer(router)

			var handler = func(channel *network.Channel) {

				callHandler = true

				payload := channel.Receive()

				options := network.NewOptions().Unmarshal(payload.Options)

				It("Should return the valid headers", func() {
					Expect(options.GetHeader("Custom-Sample")).To(Equal("Sample"))
					Expect(options.GetHeader("HP-Sample")).To(Equal("Sample"))
					Expect(options.GetHeader("GR-Sample")).To(BeEmpty())
					Expect(options.GetHeader("MQ-Sample")).To(BeEmpty())
				})

				replyOptions := network.NewOptions()

				replyOptions.SetHeader("Custom-Sample", "sample")
				replyOptions.SetHeader("HP-Sample", "sample")
				replyOptions.SetHeader("GR-Sample", "sample")
				replyOptions.SetHeader("MQ-Sample", "sample")

				channel.Send(network.BuildPayload([]byte("content reply"), replyOptions.Marshal()))

			}

			server.Listen(nethttp.SetConf("foo", "POST"), handler, nil)

			It("Should call the handler", func() {
				Expect(callHandler).To(BeTrue())
			})
		})

		Context("Listen function", func() {

			callHandler := false
			callMiddleware := false

			router := &mhttp.MockRouter{}

			router.HandleFuncMock = func(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {

				w := &mhttp.MockResponseWriter{
					Headers: make(map[string][]string),
				}

				r, _ := http.NewRequest("POST", "anything", bytes.NewReader([]byte("content")))

				r.Header = make(map[string][]string)
				r.Header["Custom-Sample"] = []string{"Sample"}
				r.Header["HP-Sample"] = []string{"Sample"}
				r.Header["GR-Sample"] = []string{"Sample"}
				r.Header["MQ-Sample"] = []string{"Sample"}

				r.Body = ioutil.NopCloser(bytes.NewReader([]byte("content")))

				f(w, r)

				It("Should contain the right headers", func() {
					Expect(w.Headers["Custom-Sample"]).To(Equal([]string{"sample"}))
					Expect(w.Headers["Hp-Sample"]).To(Equal([]string{"sample"}))
					Expect(w.Headers["Gr-Sample"]).To(BeEmpty())
					Expect(w.Headers["Mq-Sample"]).To(BeEmpty())
					Expect(string(w.Body)).To(Equal("content reply"))
					Expect(w.Status).To(Equal(201))
				})

				return &mux.Route{}
			}

			server := nethttp.NewServer(router)

			var handler = func(channel *network.Channel) {

				callHandler = true

				payload := channel.Receive()

				options := network.NewOptions().Unmarshal(payload.Options)

				It("Should return the valid headers", func() {
					Expect(options.GetHeader("Custom-Sample")).To(Equal("sample middleware"))
					Expect(options.GetHeader("HP-Sample")).To(Equal("sample middleware"))
				})
				replyOptions := network.NewOptions()

				replyOptions.SetHeader("Custom-Sample", "sample")
				replyOptions.SetHeader("HP-Sample", "sample")
				replyOptions.SetHeader("GR-Sample", "sample")
				replyOptions.SetHeader("MQ-Sample", "sample")
				replyOptions.SetHeader("Status", "201")

				channel.Send(network.BuildPayload([]byte("content reply"), replyOptions.Marshal()))

			}

			var middleware = func(handler network.Handler, channel *network.Channel) {

				callMiddleware = true

				payload := channel.Receive()

				options := network.NewOptions().Unmarshal(payload.Options)

				It("Should return the valid headers", func() {
					Expect(options.GetHeader("Custom-Sample")).To(Equal("Sample"))
					Expect(options.GetHeader("HP-Sample")).To(Equal("Sample"))
					Expect(options.GetHeader("GR-Sample")).To(BeEmpty())
					Expect(options.GetHeader("MQ-Sample")).To(BeEmpty())
				})

				replyOptions := network.NewOptions()

				replyOptions.SetHeader("Custom-Sample", "sample middleware")
				replyOptions.SetHeader("HP-Sample", "sample middleware")

				channel.Send(network.BuildPayload([]byte("content reply middleware"), replyOptions.Marshal()))

				handler(channel)
			}

			server.Listen(nethttp.SetConf("foo", "POST"), handler, middleware)

			It("Should call the handler", func() {
				Expect(callHandler).To(BeTrue())
			})

			It("Should call the handler", func() {
				Expect(callMiddleware).To(BeTrue())
			})
		})

	})

})
