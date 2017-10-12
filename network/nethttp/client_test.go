package nethttp_test

import (
	"github.com/cpapidas/pegasus/network/nethttp"

	"bytes"
	"errors"
	"github.com/cpapidas/pegasus/network"
	"github.com/cpapidas/pegasus/tests/mocks/mock_http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

var _ = Describe("Client", func() {

	Describe("Client struct", func() {

		Context("Constructor", func() {

			client := nethttp.NewClient(nil)

			It("Should not be nil", func() {
				Expect(client).ToNot(BeNil())
			})

			It("Should be type of *nethttp.Client", func() {
				Expect(reflect.ValueOf(client).String()).To(Equal("<*nethttp.Client Value>"))
			})

		})

		Context("Send function", func() {

			called := false

			mockHTTPClient := &mock_http.MockHTTPClient{}

			mockHTTPClient.DoMock = func(req *http.Request) (*http.Response, error) {

				called = true

				It("Should have the right parameters", func() {
					Expect(req.Header["Custom-Sample"]).To(Equal([]string{"sample"}))
					Expect(req.Header["Hp-Sample"]).To(Equal([]string{"sample"}))
					Expect(req.Header["Mq-Sample"]).To(BeNil())
					Expect(req.Header["Gr-Sample"]).To(BeNil())

					b, err := ioutil.ReadAll(req.Body)

					Expect(err).To(BeNil())
					Expect(string(b)).To(Equal("content"))
				})

				response := &http.Response{}
				response.Header = make(map[string][]string)
				response.Header["Custom-Sample-Reply"] = []string{"sample reply"}
				response.Header["HP-Sample-Reply"] = []string{"sample reply"}
				response.Header["MQ-Sample-Reply"] = []string{"sample reply"}
				response.Header["GR-Sample-Reply"] = []string{"sample reply"}

				response.Body = ioutil.NopCloser(bytes.NewReader([]byte("content reply")))

				return response, nil
			}

			requestOptions := network.NewOptions()
			requestOptions.SetHeader("Custom-Sample", "sample")
			requestOptions.SetHeader("HP-Sample", "sample")
			requestOptions.SetHeader("MQ-Sample", "sample")
			requestOptions.SetHeader("GR-Sample", "sample")

			client := nethttp.NewClient(mockHTTPClient)

			payload := network.BuildPayload([]byte("content"), requestOptions.Marshal())

			reply, err := client.Send(nethttp.SetConf("/whatever", "POST"), payload)

			replyOptions := network.NewOptions().Unmarshal(reply.Options)

			It("Should return a nil error", func() {
				Expect(err).To(BeNil())
			})

			It("Should call the Do function", func() {
				Expect(called).To(BeTrue())
			})

			It("Should return valid parameters", func() {
				Expect(replyOptions.GetHeader("Custom-Sample-Reply")).To(Equal("sample reply"))
				Expect(replyOptions.GetHeader("HP-Sample-Reply")).To(Equal("sample reply"))
				Expect(replyOptions.GetHeader("MQ-Sample-Reply")).To(BeEmpty())
				Expect(replyOptions.GetHeader("GR-Sample-Reply")).To(BeEmpty())
				Expect(string(reply.Body)).To(Equal("content reply"))
			})

		})

		Context("Send function request failure", func() {

			nethttp.NewRequest = func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("error")
			}

			client := nethttp.NewClient(nil)

			payload, err := client.Send(nethttp.SetConf("what", "ever"), network.Payload{})

			It("Should return nil payload", func() {
				Expect(payload).To(BeNil())
			})

			It("Should return an error", func() {
				Expect(err).ToNot(BeNil())
			})

			nethttp.NewRequest = http.NewRequest
		})

	})

	Context("Send function Do function failure", func() {

		called := false

		mockHTTPClient := &mock_http.MockHTTPClient{}

		mockHTTPClient.DoMock = func(req *http.Request) (*http.Response, error) {
			called = true
			return nil, errors.New("error")
		}

		client := nethttp.NewClient(mockHTTPClient)

		reply, err := client.Send(nethttp.SetConf("/whatever", "POST"), network.Payload{})

		It("Should call the Do function", func() {
			Expect(called).To(BeTrue())
		})

		It("Should return a nil reply", func() {
			Expect(reply).To(BeNil())
		})

		It("Should an error", func() {
			Expect(err).ToNot(BeNil())
		})

	})

	Context("Send function ReadAll function failure", func() {

		called := false

		nethttp.ReadAll = func(r io.Reader) ([]byte, error) {
			return nil, errors.New("error")
		}

		mockHTTPClient := &mock_http.MockHTTPClient{}

		mockHTTPClient.DoMock = func(req *http.Request) (*http.Response, error) {
			called = true
			response := &http.Response{}
			response.Body = ioutil.NopCloser(bytes.NewReader([]byte("content reply")))
			return response, nil
		}

		client := nethttp.NewClient(mockHTTPClient)

		reply, err := client.Send(nethttp.SetConf("/whatever", "POST"), network.Payload{})

		It("Should call the Do function", func() {
			Expect(called).To(BeTrue())
		})

		It("Should return a nil reply", func() {
			Expect(reply).To(BeNil())
		})

		It("Should an error", func() {
			Expect(err).ToNot(BeNil())
		})

		nethttp.ReadAll = ioutil.ReadAll

	})

})
