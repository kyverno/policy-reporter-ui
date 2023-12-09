package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"go.uber.org/zap"
)

type (
	DirectorOption = func(target *url.URL, req *http.Request)
	ProxyOption    = func(proxy *httputil.ReverseProxy)
)

func WithLogging() DirectorOption {
	return func(target *url.URL, req *http.Request) {
		zap.L().Debug(
			"Proxy",
			zap.String("proto", req.Proto),
			zap.String("method", req.Method),
			zap.String("referer", req.Header.Get("Referer")),
			zap.String("user-agent", req.Header.Get("User-Agent")),
			zap.String("forward-host", req.Host),
			zap.String("origin-host", target.Host),
			zap.String("path", req.URL.Path),
		)
	}
}

func WithHostOverwrite() DirectorOption {
	return func(target *url.URL, req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", target.Host)
		req.Host = target.Host
	}
}

func WithAuth(username, password string) DirectorOption {
	return func(_ *url.URL, req *http.Request) {
		req.SetBasicAuth(username, password)
	}
}

func WithCertificate(certificatePath string) ProxyOption {
	return func(proxy *httputil.ReverseProxy) {
		caCert, err := os.ReadFile(certificatePath)
		if err != nil {
			zap.L().Error("failed to read certificate", zap.Error(err), zap.String("path", certificatePath))
			return
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		proxy.Transport.(*http.Transport).TLSClientConfig.RootCAs = caCertPool
	}
}

func WithSkipTLS() ProxyOption {
	return func(proxy *httputil.ReverseProxy) {
		proxy.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = true
	}
}

func New(target *url.URL, options []DirectorOption, proxyOptions []ProxyOption) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	original := proxy.Director

	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		zap.L().Error(
			"Proxy request failed",
			zap.String("proto", req.Proto),
			zap.String("method", req.Method),
			zap.String("forward-host", req.Host),
			zap.String("origin-host", target.Host),
			zap.String("path", req.URL.Path),
			zap.Error(err),
		)
	}

	proxy.Director = func(req *http.Request) {
		for _, o := range options {
			o(target, req)
		}

		original(req)
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{},
	}

	proxy.Transport = transport

	for _, p := range proxyOptions {
		p(proxy)
	}

	return proxy
}
