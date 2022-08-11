package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func New(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	original := proxy.Director

	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", target.Host)
		req.Host = target.Host

		original(req)
	}

	proxy.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).Dial,
	}

	return proxy
}
