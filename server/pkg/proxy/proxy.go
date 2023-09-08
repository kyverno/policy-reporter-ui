package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"go.uber.org/zap"
)

func New(target *url.URL, certificatePath string, skipTLS, overwriteHost, logging bool, username, password string) *httputil.ReverseProxy {
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
		if logging {
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

		if overwriteHost {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", target.Host)
			req.Host = target.Host
		}

		if username != "" && password != "" {
			req.SetBasicAuth(username, password)
		}

		original(req)
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTLS,
		},
	}

	proxy.Transport = transport

	if certificatePath != "" {
		caCert, err := ioutil.ReadFile(certificatePath)
		if err != nil {
			log.Printf("[ERROR] failed to read certificate: %s\n", certificatePath)
			return proxy
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		transport.TLSClientConfig.RootCAs = caCertPool
	}

	return proxy
}
