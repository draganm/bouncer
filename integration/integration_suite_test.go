package integration_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/draganm/web-interceptor/proxy"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var upgrader = websocket.Upgrader{} // use default options

var _ = BeforeSuite(func(done Done) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
		// w.WriteHeader(200)
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
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

		Context("When upgrading to an WebSockets connection", func() {
			var err error
			var c *websocket.Conn
			BeforeEach(func() {
				c, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
				Expect(err).ToNot(HaveOccurred())
			})

			It("Should not fail", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			Context("When I send a text message through the socket", func() {
				BeforeEach(func() {
					Expect(c.WriteMessage(websocket.TextMessage, []byte("test"))).To(Succeed())
				})
				Context("And I read the message sent back", func() {
					var tpe int
					var message []byte
					var err error
					BeforeEach(func() {
						tpe, message, err = c.ReadMessage()
					})

					It("Should not to fail", func() {
						Expect(err).ToNot(HaveOccurred())
					})
					It("Should have text message type", func() {
						Expect(tpe).To(Equal(websocket.TextMessage))
					})
					It("Should have the same message body as sent", func() {
						Expect(string(message)).To(Equal("test"))
					})
				})
			})

		})

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
