package cluster

import (
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/auth"
	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
)

type BasicAuth struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

type AccessControl struct {
	Emails []string `koanf:"emails"`
	Groups []string `koanf:"groups"`
}

type Plugin struct {
	Name        string    `koanf:"name"`
	Host        string    `koanf:"host"`
	HTTP2       bool      `koanf:"http2"`
	SkipTLS     bool      `koanf:"skipTLS"`
	Certificate string    `koanf:"certificate"`
	SecretRef   string    `koanf:"secretRef"`
	BasicAuth   BasicAuth `koanf:"basicAuth"`
	Logging     bool      `koanf:"-"`
}

func (a Plugin) FromValues(values secrets.Values) Plugin {
	if values.Host != "" {
		a.Host = values.Host
	}
	if values.Certificate != "" {
		a.Certificate = values.Certificate
	}
	if values.SkipTLS {
		a.SkipTLS = values.SkipTLS
	}
	if values.Username != "" {
		a.BasicAuth.Username = values.Username
	}
	if values.Password != "" {
		a.BasicAuth.Password = values.Password
	}

	return a
}

type Config struct {
	Name          string        `koanf:"name"`
	Host          string        `koanf:"host"`
	HTTP2         bool          `koanf:"http2"`
	Plugins       []Plugin      `koanf:"plugins"`
	SkipTLS       bool          `koanf:"skipTLS"`
	Certificate   string        `koanf:"certificate"`
	SecretRef     string        `koanf:"secretRef"`
	BasicAuth     BasicAuth     `koanf:"basicAuth"`
	AccessControl AccessControl `koanf:"accessControl"`
	Logging       bool          `koanf:"-"`
}

func (a Config) FromValues(values secrets.Values) Config {
	if values.Host != "" {
		a.Host = values.Host
	}
	if values.Certificate != "" {
		a.Certificate = values.Certificate
	}
	if values.SkipTLS {
		a.SkipTLS = values.SkipTLS
	}
	if values.SecretRef != "" {
		a.SecretRef = values.SecretRef
	}
	if values.Username != "" {
		a.BasicAuth.Username = values.Username
	}
	if values.Password != "" {
		a.BasicAuth.Password = values.Password
	}

	for _, p := range values.Plugins {
		a.Plugins = append(a.Plugins, Plugin{
			Host:        p.Host,
			Name:        p.Name,
			SkipTLS:     p.SkipTLS,
			Certificate: p.Certificate,
			BasicAuth: BasicAuth{
				Username: p.Username,
				Password: p.Password,
			},
		})
	}

	return a
}

type Cluster struct {
	Name    string
	Core    *core.Client
	Plugins map[string]*plugin.Client
	auth.Permissions
}
