package config

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ctrserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/kyverno/policy-reporter-ui/pkg/server"
)

func TestTLSConfigLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`server:
  tls:
    enabled: true
    certFile: /ui/tls.crt
    keyFile: /ui/tls.key
    clientAuth:
      enabled: true
      caFile: /ui/ca.crt
  metrics:
    enabled: true
    tls:
      enabled: true
      certFile: /metrics/tls.crt
      keyFile: /metrics/tls.key
`), 0o600); err != nil {
		t.Fatal(err)
	}
	c := &Config{}
	if err := Load(c, path); err != nil {
		t.Fatal(err)
	}
	if !c.Server.TLS.Enabled || !c.Server.TLS.ClientAuth.Enabled || c.Server.TLS.CertFile != "/ui/tls.crt" || c.Server.TLS.ClientAuth.CAFile != "/ui/ca.crt" {
		t.Fatalf("UI TLS config not decoded: %+v", c.Server.TLS)
	}
	if !c.Server.Metrics.TLS.Enabled || c.Server.Metrics.TLS.ClientAuth.Enabled || c.Server.Metrics.TLS.CertFile != "/metrics/tls.crt" {
		t.Fatalf("metrics TLS config not independent: %+v", c.Server.Metrics.TLS)
	}
	t.Setenv("PR_SERVER_TLS_ENABLED", "true")
	t.Setenv("PR_SERVER_TLS_CERTFILE", "/env/tls.crt")
	t.Setenv("PR_SERVER_TLS_KEYFILE", "/env/tls.key")
	t.Setenv("PR_SERVER_TLS_CLIENTAUTH_ENABLED", "true")
	t.Setenv("PR_SERVER_TLS_CLIENTAUTH_CAFILE", "/env/ca.crt")
	t.Setenv("PR_SERVER_METRICS_TLS_ENABLED", "true")
	c = &Config{}
	t.Chdir(t.TempDir())
	if err := Load(c, ""); err != nil {
		t.Fatal(err)
	}
	if !c.Server.TLS.Enabled || !c.Server.TLS.ClientAuth.Enabled || c.Server.TLS.CertFile != "/env/tls.crt" || c.Server.TLS.ClientAuth.CAFile != "/env/ca.crt" || !c.Server.Metrics.TLS.Enabled {
		t.Fatalf("environment TLS config not decoded: %+v", c.Server)
	}
}

func TestMetricsTLSOptions(t *testing.T) {
	t.Parallel()
	c := &Config{}
	options, err := NewResolver(c).metricsOptions()
	if err != nil || options.SecureServing || options.BindAddress != "0" {
		t.Fatalf("default metrics: %+v %v", options, err)
	}
	c.Server.Metrics.Enabled = true
	c.Server.Metrics.Port = 9090
	options, err = NewResolver(c).metricsOptions()
	if err != nil || options.SecureServing || options.BindAddress != ":9090" {
		t.Fatalf("HTTP metrics: %+v %v", options, err)
	}
	fixture := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer fixture.Close()
	cert := fixture.TLS.Certificates[0]
	key, err := x509.MarshalPKCS8PrivateKey(cert.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	certFile := filepath.Join(dir, "tls.crt")
	keyFile := filepath.Join(dir, "tls.key")
	if err := os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate[0]}), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: key}), 0o600); err != nil {
		t.Fatal(err)
	}
	c.Server.Metrics.TLS = server.TLSOptions{Enabled: true, CertFile: certFile, KeyFile: keyFile, ClientAuth: server.ClientAuthOptions{Enabled: true, CAFile: certFile}}
	resolver := NewResolver(c)
	options, err = resolver.metricsOptions()
	if err != nil {
		t.Fatal(err)
	}
	cfg := &tls.Config{}
	for _, apply := range options.TLSOpts {
		apply(cfg)
	}
	if !options.SecureServing || cfg.ClientAuth != tls.RequireAndVerifyClientCert || cfg.GetCertificate == nil || cfg.ClientCAs == nil || cfg.MinVersion != tls.VersionTLS12 {
		t.Fatalf("metrics TLS lost settings: %+v", cfg)
	}
	loaded, err := cfg.GetCertificate(nil)
	if err != nil || string(loaded.Certificate[0]) != string(cert.Certificate[0]) {
		t.Fatal("metrics serving certificate replaced")
	}
	c.Server.TLS = c.Server.Metrics.TLS
	c.Server.TLS.ClientAuth.Mode = "oidc-or-certificate"
	c.Server.TLS.ClientAuth.Identities = []server.CertificateIdentity{{Fingerprint: strings.Repeat("0", 64), Email: "client@example.test"}}
	if _, err := NewResolver(c).Server(context.Background()); err == nil || !strings.Contains(err.Error(), "requires OpenID Connect") {
		t.Fatalf("certificate login accepted without OIDC: %v", err)
	}
	c.Server.TLS = server.TLSOptions{}
	c.Server.Metrics.TLS.ClientAuth.Mode = "oidc-or-certificate"
	if _, err := NewResolver(c).metricsOptions(); err == nil {
		t.Fatal("metrics accepted OIDC mode")
	}
	c.Server.Metrics.TLS.ClientAuth.Mode = ""
	if err := os.Remove(certFile); err != nil {
		t.Fatal(err)
	}
	if _, err := resolver.metricsOptions(); err != nil {
		t.Fatalf("certificate reloaded after preparation: %v", err)
	}
	if c.Server.TLS.Enabled {
		t.Fatal("metrics changed UI TLS")
	}
	c.Server.Metrics.Enabled = false
	if _, err := NewResolver(c).metricsOptions(); err == nil {
		t.Fatal("enabled TLS on disabled metrics accepted")
	}
}

func TestMetricsListener(t *testing.T) {
	t.Parallel()
	fixture := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	key := fixture.TLS.Certificates[0].PrivateKey.(crypto.Signer)
	fixture.Close()
	template := &x509.Certificate{SerialNumber: big.NewInt(1), NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}, IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}}
	der, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	certFile := filepath.Join(dir, "tls.crt")
	keyFile := filepath.Join(dir, "tls.key")
	if err := os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}), 0o600); err != nil {
		t.Fatal(err)
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	parsed, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	roots.AddCert(parsed)
	for _, mode := range []string{"HTTP", "TLS", "mTLS"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			c := &Config{}
			c.Server.Metrics.Enabled = true
			if mode != "HTTP" {
				c.Server.Metrics.TLS = server.TLSOptions{Enabled: true, CertFile: certFile, KeyFile: keyFile}
			}
			if mode == "mTLS" {
				c.Server.Metrics.TLS.ClientAuth = server.ClientAuthOptions{Enabled: true, CAFile: certFile}
			}
			options, err := NewResolver(c).metricsOptions()
			if err != nil {
				t.Fatal(err)
			}
			options.BindAddress = "127.0.0.1:0"
			listener, err := ctrserver.NewServer(options, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- listener.Start(ctx) }()
			defer func() {
				cancel()
				if err := <-done; err != nil {
					t.Errorf("metrics stopped: %v", err)
				}
			}()
			address := listener.(interface{ GetBindAddr() string })
			deadline := time.Now().Add(5 * time.Second)
			for address.GetBindAddr() == "" {
				if time.Now().After(deadline) {
					t.Fatal("metrics did not start")
				}
				time.Sleep(time.Millisecond)
			}
			scheme := "http"
			if mode != "HTTP" {
				scheme = "https"
			}
			for _, withClient := range []bool{false, true} {
				cfg := &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
				if withClient {
					cfg.Certificates = []tls.Certificate{cert}
				}
				transport := &http.Transport{TLSClientConfig: cfg}
				defer transport.CloseIdleConnections()
				client := &http.Client{Transport: transport, Timeout: 3 * time.Second}
				req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, scheme+"://"+address.GetBindAddr()+"/metrics", nil)
				if err != nil {
					t.Fatal(err)
				}
				resp, err := client.Do(req)
				want := mode != "mTLS" || withClient
				if (err == nil) != want {
					t.Fatalf("client=%v success=%v want %v: %v", withClient, err == nil, want, err)
				}
				if err == nil {
					if _, err := io.Copy(io.Discard, resp.Body); err != nil {
						resp.Body.Close()
						t.Fatal(err)
					}
					resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						t.Fatalf("status %d", resp.StatusCode)
					}
				}
			}
		})
	}
}
