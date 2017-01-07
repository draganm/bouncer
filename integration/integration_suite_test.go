package integration_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/draganm/web-interceptor/proxy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func(done Done) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
		w.WriteHeader(200)
	})
	go http.ListenAndServe("localhost:8081", mux)
	for {
		_, err := http.Get("http://localhost:8081")
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(time.Millisecond * 10)
	}
	go proxy.Proxy(":8080", "http://localhost:8081")
	for {
		_, err := http.Get("http://localhost:8080")
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}

	close(done)
}, 3.0)

var _ = Describe("simple proxy", func() {

	Describe("GET request", func() {
		Context("When requesting existing resource through the proxy", func() {
			var response *http.Response
			BeforeEach(func() {
				var err error
				response, err = http.Get("http://localhost:8080/")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				response.Body.Close()
			})

			It("Should respond with 200 status code", func() {
				Expect(response.StatusCode).To(Equal(200))
			})

			It("Should forward original response body", func() {
				data, err := ioutil.ReadAll(response.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(Equal("test"))
			})

		})
	})
})
