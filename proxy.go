package bouncer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/koding/websocketproxy"
	"github.com/urfave/negroni"
)

func Proxy(addr, remote string, handlers ...negroni.Handler) error {
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

	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Header.Get("Connection")) == "upgrade" {
			wsProxy.ServeHTTP(w, r)
			return
		}
		proxy.ServeHTTP(w, r)
	})

	n := negroni.New(handlers...)
	n.UseHandler(m)
	n.Run(addr)
	return nil
	// return http.ListenAndServe(addr, m)
}
