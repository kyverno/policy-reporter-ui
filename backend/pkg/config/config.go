package config

import (
	"strings"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"github.com/kyverno/policy-reporter-ui/pkg/logging"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

// BasicAuth configuration
type BasicAuth struct {
	Username  string `koanf:"username"`
	Password  string `koanf:"password"`
	SecretRef string `koanf:"secretRef"`
}

type OpenIDConnect struct {
	Enabled      bool     `koanf:"enabled"`
	SecretRef    string   `koanf:"secretRef"`
	DiscoveryURL string   `koanf:"discoveryUrl"`
	CallbackURL  string   `koanf:"callbackUrl"`
	ClientID     string   `koanf:"clientId"`
	ClientSecret string   `koanf:"clientSecret"`
	GroupClaim   string   `koanf:"groupClaim"`
	Scopes       []string `koanf:"scopes"`
}

func (a OpenIDConnect) BasePath() string {
	return utils.BasePath(a.CallbackURL)
}

func (a OpenIDConnect) Callback() string {
	if strings.HasSuffix(a.CallbackURL, "/callback") {
		return a.CallbackURL
	}

	return strings.TrimSuffix(a.CallbackURL, "/") + "/callback"
}

func (a OpenIDConnect) Discovery() string {
	if strings.HasSuffix(a.DiscoveryURL, "/.well-known/openid-configuration") {
		return a.DiscoveryURL
	}

	return strings.TrimSuffix(a.DiscoveryURL, "/") + "/.well-known/openid-configuration"
}

func (a OpenIDConnect) FromValues(values secrets.Values) OpenIDConnect {
	if values.DiscoveryURL != "" {
		a.DiscoveryURL = values.DiscoveryURL
	}
	if values.ClientID != "" {
		a.ClientID = values.ClientID
	}
	if values.ClientSecret != "" {
		a.ClientSecret = values.ClientSecret
	}

	return a
}

type OAuth struct {
	Enabled      bool     `koanf:"enabled"`
	SecretRef    string   `koanf:"secretRef"`
	Provider     string   `koanf:"provider"`
	CallbackURL  string   `koanf:"callbackUrl"`
	ClientID     string   `koanf:"clientId"`
	ClientSecret string   `koanf:"clientSecret"`
	Scopes       []string `koanf:"scopes"`
}

func (a OAuth) FromValues(values secrets.Values) OAuth {
	if values.Provider != "" {
		a.Provider = values.Provider
	}
	if values.ClientSecret != "" {
		a.ClientSecret = values.ClientSecret
	}
	if values.ClientID != "" {
		a.ClientID = values.ClientID
	}

	return a
}

func (a OAuth) BasePath() string {
	return utils.BasePath(a.CallbackURL)
}

type Plugin struct {
	Name        string    `koanf:"name"`
	Host        string    `koanf:"host"`
	SkipTLS     bool      `koanf:"skipTLS"`
	Certificate string    `koanf:"certificate"`
	SecretRef   string    `koanf:"secretRef"`
	BasicAuth   BasicAuth `koanf:"basicAuth"`
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

// Cluster configuration
type Cluster struct {
	Name          string        `koanf:"name"`
	Host          string        `koanf:"host"`
	Plugins       []Plugin      `koanf:"plugins"`
	SkipTLS       bool          `koanf:"skipTLS"`
	Certificate   string        `koanf:"certificate"`
	SecretRef     string        `koanf:"secretRef"`
	BasicAuth     BasicAuth     `koanf:"basicAuth"`
	AccessControl AccessControl `koanf:"accessControl"`
}

func (a Cluster) FromValues(values secrets.Values) Cluster {
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

type Logo struct {
	Path     string `koanf:"path"`
	Disabled bool   `koanf:"disabled"`
}

type UI struct {
	Path        string `koanf:"path"`
	Banner      string `koanf:"banner"`
	Disabled    bool   `koanf:"enabled"`
	DisplayMode string `koanf:"displayMode"`
	Logo        Logo   `koanf:"logo"`
}

type Server struct {
	Port          int  `koanf:"port"`
	CORS          bool `koanf:"cors"`
	Debug         bool `koanf:"debug"`
	OverwriteHost bool `koanf:"overwriteHost"`
}

type Source struct {
	Name       string `koanf:"name"`
	Exceptions bool   `koanf:"exceptions"`
	ViewType   string `koanf:"type"`
	Excludes   struct {
		NamespaceKinds []string `koanf:"namespaceKinds"`
		ClusterKinds   []string `koanf:"clusterKinds"`
		Results        []string `koanf:"results"`
		Severities     []string `koanf:"severities"`
	} `koanf:"excludes"`
}

type AccessControl struct {
	Emails []string `koanf:"emails"`
	Groups []string `koanf:"groups"`
}

type Boards struct {
	AccessControl AccessControl `koanf:"accessControl"`
}

type CustomBoard struct {
	Name          string        `koanf:"name"`
	AccessControl AccessControl `koanf:"accessControl"`
	Namespaces    struct {
		Selector map[string]string `koanf:"selector"`
		List     []string          `koanf:"list"`
	} `koanf:"namespaces"`
	Sources struct {
		List []string `koanf:"list"`
	} `koanf:"sources"`
	PolicyReports struct {
		Selector map[string]string `koanf:"selector"`
	} `koanf:"policyReports"`
	ClusterScope struct {
		Enabled bool `koanf:"enabled"`
	} `koanf:"clusterScope"`
}

// Config structure
type Config struct {
	KubeConfig    clientcmd.ConfigOverrides
	Namespace     string         `koanf:"namespace"`
	TempDir       string         `koanf:"tempDir"`
	Clusters      []Cluster      `koanf:"clusters"`
	Sources       []Source       `koanf:"sources"`
	Server        Server         `koanf:"server"`
	UI            UI             `koanf:"ui"`
	Logging       logging.Config `koanf:"logging"`
	OpenIDConnect OpenIDConnect  `koanf:"openIDConnect"`
	OAuth         OAuth          `koanf:"oauth"`
	Boards        Boards         `mastructure:"boards"`
	CustomBoards  []CustomBoard  `koanf:"customBoards"`
	Local         bool           `koanf:"local"`
}

func (c *Config) AuthEnabled() bool {
	return c.OAuth.Enabled || c.OpenIDConnect.Enabled
}

func (c *Config) AuthBasePath() string {
	if c.OAuth.Enabled {
		return c.OAuth.BasePath()
	}

	return c.OpenIDConnect.BasePath()
}

func (c *Config) AuthGroupClaim() string {
	if c.OAuth.Enabled {
		return ""
	}

	return c.OpenIDConnect.GroupClaim
}
