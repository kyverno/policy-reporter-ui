package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"go.uber.org/zap"

	"github.com/kyverno/policy-reporter-ui/pkg/api"
)

type (
	DirectorOption = func(target *url.URL, req *httputil.ProxyRequest)
	ProxyOption    = func(proxy *httputil.ReverseProxy)
)

func WithLogging() DirectorOption {
	return func(target *url.URL, req *httputil.ProxyRequest) {
		zap.L().Debug(
			"ProxyRequest.Int",
			zap.String("proto", req.In.Proto),
			zap.String("method", req.In.Method),
			zap.String("referer", req.In.Header.Get("Referer")),
			zap.String("user-agent", req.In.Header.Get("User-Agent")),
			zap.String("forward-host", req.In.Host),
			zap.String("origin-host", target.Host),
			zap.String("path", req.In.URL.Path),
		)
		zap.L().Debug(
			"ProxyRequest.Out",
			zap.String("proto", req.Out.Proto),
			zap.String("method", req.Out.Method),
			zap.String("referer", req.Out.Header.Get("Referer")),
			zap.String("user-agent", req.Out.Header.Get("User-Agent")),
			zap.String("forward-host", req.Out.Host),
			zap.String("origin-host", target.Host),
			zap.String("path", req.Out.URL.Path),
		)
	}
}

func WithHostOverwrite() DirectorOption {
	return func(target *url.URL, req *httputil.ProxyRequest) {
		req.Out.Header.Add("X-Forwarded-Host", req.In.Host)
		req.Out.Header.Add("X-Origin-Host", target.Host)
		req.Out.Host = target.Host
	}
}

func WithAuth(username, password string) DirectorOption {
	return func(_ *url.URL, req *httputil.ProxyRequest) {
		req.Out.SetBasicAuth(username, password)
	}
}

func WithCertificate(certificatePath string) ProxyOption {
	return func(proxy *httputil.ReverseProxy) {
		pool, err := api.LoadCerts(certificatePath)
		if err != nil {
			zap.L().Error("failed to read certificate", zap.Error(err), zap.String("path", certificatePath))
			return
		}

		proxy.Transport.(*http.Transport).TLSClientConfig.RootCAs = pool
	}
}

func WithSkipTLS() ProxyOption {
	return func(proxy *httputil.ReverseProxy) {
		proxy.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = true
	}
}

func New(target *url.URL, options []DirectorOption, proxyOptions []ProxyOption) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	original := proxy.Rewrite

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

	proxy.Rewrite = func(req *httputil.ProxyRequest) {
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
		TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}

	proxy.Transport = transport

	for _, p := range proxyOptions {
		p(proxy)
	}

	return proxy
}
