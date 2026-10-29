package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTLSOptions(t *testing.T) {
	ca, caKey := testCertificate(t, nil, nil, true, x509.ExtKeyUsageAny, false)
	cert, key := testCertificate(t, ca, caKey, false, x509.ExtKeyUsageServerAuth, false)
	dir := t.TempDir()
	write := func(name, kind string, data []byte) string {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: kind, Bytes: data}), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	certFile := write("tls.crt", "CERTIFICATE", cert.Raw)
	keyFile := write("tls.key", "EC PRIVATE KEY", keyDER)
	caFile := write("ca.crt", "CERTIFICATE", ca.Raw)
	_, otherKey := testCertificate(t, nil, nil, true, x509.ExtKeyUsageAny, false)
	otherDER, err := x509.MarshalECPrivateKey(otherKey)
	if err != nil {
		t.Fatal(err)
	}
	otherKeyFile := write("other.key", "EC PRIVATE KEY", otherDER)
	for name, options := range map[string]TLSOptions{
		"disabled TLS with certificate": {CertFile: certFile},
		"mTLS without TLS":              {ClientAuth: ClientAuthOptions{Enabled: true}},
		"unknown mode":                  {Enabled: true, CertFile: certFile, KeyFile: keyFile, ClientAuth: ClientAuthOptions{Enabled: true, CAFile: caFile, Mode: "unknown"}},
		"mapping without identities":    {Enabled: true, CertFile: certFile, KeyFile: keyFile, ClientAuth: ClientAuthOptions{Enabled: true, CAFile: caFile, Mode: "oidc-or-certificate"}},
		"invalid fingerprint":           {Enabled: true, CertFile: certFile, KeyFile: keyFile, ClientAuth: ClientAuthOptions{Enabled: true, CAFile: caFile, Mode: "oidc-or-certificate", Identities: []CertificateIdentity{{Fingerprint: "bad", Email: "user@example.test"}}}},
		"mismatched key":                {Enabled: true, CertFile: certFile, KeyFile: otherKeyFile},
		"missing key":                   {Enabled: true, CertFile: certFile},
		"unreadable certificate":        {Enabled: true, CertFile: "missing", KeyFile: keyFile},
		"malformed certificate":         {Enabled: true, CertFile: keyFile, KeyFile: keyFile},
		"missing CA":                    {Enabled: true, CertFile: certFile, KeyFile: keyFile, ClientAuth: ClientAuthOptions{Enabled: true}},
		"unreadable CA":                 {Enabled: true, CertFile: certFile, KeyFile: keyFile, ClientAuth: ClientAuthOptions{Enabled: true, CAFile: "missing"}},
		"malformed CA":                  {Enabled: true, CertFile: certFile, KeyFile: keyFile, ClientAuth: ClientAuthOptions{Enabled: true, CAFile: keyFile}},
		"disabled mTLS with CA":         {Enabled: true, CertFile: certFile, KeyFile: keyFile, ClientAuth: ClientAuthOptions{CAFile: caFile}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := options.Config(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
	if cfg, err := (TLSOptions{}).Config(); err != nil || cfg != nil {
		t.Fatalf("default must keep HTTP: %v %v", cfg, err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(ca)
	otherCA, otherKey := testCertificate(t, nil, nil, true, x509.ExtKeyUsageAny, false)
	for _, mode := range []string{"TLS", "mTLS", "OIDC-or-certificate"} {
		mutual := mode == "mTLS"
		verifyClient := mode != "TLS"
		options := TLSOptions{Enabled: true, CertFile: certFile, KeyFile: keyFile, ClientAuth: ClientAuthOptions{Enabled: verifyClient}}
		if verifyClient {
			options.ClientAuth.CAFile = caFile
		}
		if mode == "OIDC-or-certificate" {
			options.ClientAuth.Mode = "oidc-or-certificate"
			options.ClientAuth.Identities = []CertificateIdentity{{Fingerprint: "0000000000000000000000000000000000000000000000000000000000000000", Email: "user@example.test"}}
		}
		cfg, err := options.Config()
		if err != nil {
			t.Fatal(err)
		}
		if cfg.MinVersion != tls.VersionTLS12 {
			t.Fatal("TLS minimum not enforced")
		}
		ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }))
		ts.Config.ErrorLog = log.New(io.Discard, "", 0)
		ts.TLS = cfg
		ts.StartTLS()
		defer ts.Close()
		for _, tc := range []struct {
			name                string
			parent              *x509.Certificate
			key                 *ecdsa.PrivateKey
			usage               x509.ExtKeyUsage
			expired, send, want bool
		}{
			{name: "no client", want: !mutual},
			{name: "trusted", parent: ca, key: caKey, usage: x509.ExtKeyUsageClientAuth, send: true, want: true},
			{name: "untrusted", parent: otherCA, key: otherKey, usage: x509.ExtKeyUsageClientAuth, send: true, want: !verifyClient},
			{name: "expired", parent: ca, key: caKey, usage: x509.ExtKeyUsageClientAuth, expired: true, send: true, want: !verifyClient},
			{name: "wrong usage", parent: ca, key: caKey, usage: x509.ExtKeyUsageServerAuth, send: true, want: !verifyClient},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				clientCfg := &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
				if tc.send {
					leaf, leafKey := testCertificate(t, tc.parent, tc.key, false, tc.usage, tc.expired)
					clientCert := tls.Certificate{Certificate: [][]byte{leaf.Raw}, PrivateKey: leafKey}
					// Present even untrusted test certificates regardless of advertised CA names.
					clientCfg.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) { return &clientCert, nil }
				}
				transport := &http.Transport{TLSClientConfig: clientCfg}
				defer transport.CloseIdleConnections()
				client := &http.Client{Transport: transport, Timeout: 3 * time.Second}
				resp, err := client.Get(ts.URL)
				if err == nil {
					defer resp.Body.Close()
				}
				if (err == nil) != tc.want {
					t.Fatalf("success=%v, want %v: %v", err == nil, tc.want, err)
				}
			})
		}
		resp, err := http.Get("http://" + ts.Listener.Addr().String())
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusNoContent {
				t.Fatal("plaintext accepted")
			}
		}
	}
}

func testCertificate(t *testing.T, parent *x509.Certificate, parentKey *ecdsa.PrivateKey, ca bool, usage x509.ExtKeyUsage, expired bool) (*x509.Certificate, *ecdsa.PrivateKey) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatal(err)
	}
	cert := &x509.Certificate{SerialNumber: serial, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), BasicConstraintsValid: true, IsCA: ca, KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{usage}, IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}}
	if ca {
		cert.KeyUsage |= x509.KeyUsageCertSign
	}
	if expired {
		cert.NotBefore = time.Now().Add(-2 * time.Hour)
		cert.NotAfter = time.Now().Add(-time.Hour)
	}
	if parent == nil {
		parent = cert
		parentKey = key
	}
	der, err := x509.CreateCertificate(rand.Reader, cert, parent, &key.PublicKey, parentKey)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return parsed, key
}
