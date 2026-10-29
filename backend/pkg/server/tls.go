package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
)

type TLSOptions struct {
	Enabled    bool              `koanf:"enabled"`
	CertFile   string            `koanf:"certFile"`
	KeyFile    string            `koanf:"keyFile"`
	ClientAuth ClientAuthOptions `koanf:"clientAuth"`
}

type CertificateIdentity struct {
	Fingerprint string   `koanf:"fingerprint"`
	Email       string   `koanf:"email"`
	Name        string   `koanf:"name"`
	Groups      []string `koanf:"groups"`
}

type ClientAuthOptions struct {
	Mode       string                `koanf:"mode"`
	Identities []CertificateIdentity `koanf:"identities"`
	Enabled    bool                  `koanf:"enabled"`
	CAFile     string                `koanf:"caFile"`
}

func (o TLSOptions) Config() (*tls.Config, error) {
	if !o.Enabled {
		if o.CertFile != "" || o.KeyFile != "" || o.ClientAuth.Enabled || o.ClientAuth.CAFile != "" || o.ClientAuth.Mode != "" || len(o.ClientAuth.Identities) != 0 {
			return nil, fmt.Errorf("server TLS options require server.tls.enabled")
		}
		return nil, nil
	}
	if o.CertFile == "" || o.KeyFile == "" {
		return nil, fmt.Errorf("server TLS requires certFile and keyFile")
	}
	if !o.ClientAuth.Enabled && (o.ClientAuth.CAFile != "" || o.ClientAuth.Mode != "" || len(o.ClientAuth.Identities) != 0) {
		return nil, fmt.Errorf("server TLS client authentication options require server.tls.clientAuth.enabled")
	}
	if o.ClientAuth.Mode != "" && o.ClientAuth.Mode != "required" && o.ClientAuth.Mode != "oidc-or-certificate" {
		return nil, fmt.Errorf("unsupported server TLS client authentication mode")
	}
	if o.ClientAuth.Mode != "oidc-or-certificate" && len(o.ClientAuth.Identities) != 0 {
		return nil, fmt.Errorf("certificate identities require oidc-or-certificate mode")
	}
	if o.ClientAuth.Mode == "oidc-or-certificate" && len(o.ClientAuth.Identities) == 0 {
		return nil, fmt.Errorf("oidc-or-certificate mode requires certificate identities")
	}
	seen := map[string]bool{}
	for _, identity := range o.ClientAuth.Identities {
		fingerprint, err := hex.DecodeString(identity.Fingerprint)
		if err != nil || len(fingerprint) != 32 || identity.Email == "" {
			return nil, fmt.Errorf("certificate identity requires a SHA-256 fingerprint and email")
		}
		canonical := hex.EncodeToString(fingerprint)
		if seen[canonical] {
			return nil, fmt.Errorf("duplicate certificate identity fingerprint")
		}
		seen[canonical] = true
	}
	cert, err := tls.LoadX509KeyPair(o.CertFile, o.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("load server TLS certificate: %w", err)
	}
	config := &tls.Config{MinVersion: tls.VersionTLS12, Certificates: []tls.Certificate{cert}}
	if o.ClientAuth.Enabled {
		if o.ClientAuth.CAFile == "" {
			return nil, fmt.Errorf("server TLS client authentication requires caFile")
		}
		data, err := os.ReadFile(o.ClientAuth.CAFile)
		if err != nil {
			return nil, fmt.Errorf("read server TLS client CA: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("server TLS client CA contains no certificates")
		}
		config.ClientCAs = pool
		config.ClientAuth = tls.RequireAndVerifyClientCert
		if o.ClientAuth.Mode == "oidc-or-certificate" {
			config.ClientAuth = tls.VerifyClientCertIfGiven
		}
	}
	return config, nil
}
