package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Proxy(addr, remote string) error {
	u, err := url.Parse(remote)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	return http.ListenAndServe(addr, nil)
}
