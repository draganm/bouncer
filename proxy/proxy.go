package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/koding/websocketproxy"
)

func Proxy(addr, remote string) error {
	u, err := url.Parse(remote)
	if err != nil {
		return err
	}
	wsURL, err := url.Parse(strings.Replace(remote, "http", "ws", 1))
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	wsProxy := websocketproxy.NewProxy(wsURL)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Header.Get("Connection")) == "upgrade" {
			wsProxy.ServeHTTP(w, r)
			return
		}
		proxy.ServeHTTP(w, r)
	})

	return http.ListenAndServe(addr, nil)
}
