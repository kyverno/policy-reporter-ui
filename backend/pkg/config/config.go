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
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	SecretRef string `mapstructure:"secretRef"`
}

type OpenIDConnect struct {
	Enabled      bool     `mapstructure:"enabled"`
	SecretRef    string   `mapstructure:"secretRef"`
	DiscoveryURL string   `mapstructure:"discoveryUrl"`
	CallbackURL  string   `mapstructure:"callbackUrl"`
	ClientID     string   `mapstructure:"clientId"`
	ClientSecret string   `mapstructure:"clientSecret"`
	Scopes       []string `mapstructure:"scopes"`
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
	Enabled      bool     `mapstructure:"enabled"`
	SecretRef    string   `mapstructure:"secretRef"`
	Provider     string   `mapstructure:"provider"`
	CallbackURL  string   `mapstructure:"callbackUrl"`
	ClientID     string   `mapstructure:"clientId"`
	ClientSecret string   `mapstructure:"clientSecret"`
	Scopes       []string `mapstructure:"scopes"`
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
	Name        string    `mapstructure:"name"`
	Host        string    `mapstructure:"host"`
	SkipTLS     bool      `mapstructure:"skipTLS"`
	Certificate string    `mapstructure:"certificate"`
	SecretRef   string    `mapstructure:"secretRef"`
	BasicAuth   BasicAuth `mapstructure:"basicAuth"`
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
	Name          string        `mapstructure:"name"`
	Host          string        `mapstructure:"host"`
	Plugins       []Plugin      `mapstructure:"plugins"`
	SkipTLS       bool          `mapstructure:"skipTLS"`
	Certificate   string        `mapstructure:"certificate"`
	SecretRef     string        `mapstructure:"secretRef"`
	BasicAuth     BasicAuth     `mapstructure:"basicAuth"`
	AccessControl AccessControl `mapstructure:"accessControl"`
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

type UI struct {
	Path        string `mapstructure:"path"`
	Banner      string `mapstructure:"banner"`
	Disabled    bool   `mapstructure:"enabled"`
	DisplayMode string `mapstructure:"displayMode"`
}

type Server struct {
	Port          int  `mapstructure:"port"`
	CORS          bool `mapstructure:"cors"`
	Debug         bool `mapstructure:"debug"`
	OverwriteHost bool `mapstructure:"overwriteHost"`
}

type Source struct {
	Name       string `mapstructure:"name"`
	Exceptions bool   `mapstructure:"exceptions"`
	ViewType   string `mapstructure:"type"`
	Excludes   struct {
		NamespaceKinds []string `mapstructure:"namespaceKinds"`
		ClusterKinds   []string `mapstructure:"clusterKinds"`
		Results        []string `mapstructure:"results"`
		Severities     []string `mapstructure:"severities"`
	} `mapstructure:"excludes"`
}

type AccessControl struct {
	Emails []string `mapstructure:"emails"`
}

type Boards struct {
	AccessControl AccessControl `mapstructure:"accessControl"`
}

type CustomBoard struct {
	Name          string        `mapstructure:"name"`
	AccessControl AccessControl `mapstructure:"accessControl"`
	Namespaces    struct {
		Selector map[string]string `mapstructure:"selector"`
		List     []string          `mapstructure:"list"`
	} `mapstructure:"namespaces"`
	Sources struct {
		List []string `mapstructure:"list"`
	} `mapstructure:"sources"`
	PolicyReports struct {
		Selector map[string]string `mapstructure:"selector"`
	} `mapstructure:"policyReports"`
	ClusterScope struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"clusterScope"`
}

// Config structure
type Config struct {
	KubeConfig    clientcmd.ConfigOverrides
	Namespace     string         `mapstructure:"namespace"`
	TempDir       string         `mapstructure:"tempDir"`
	Clusters      []Cluster      `mapstructure:"clusters"`
	Sources       []Source       `mapstructure:"sources"`
	Server        Server         `mapstructure:"server"`
	UI            UI             `mapstructure:"ui"`
	Logging       logging.Config `mapstructure:"logging"`
	OpenIDConnect OpenIDConnect  `mapstructure:"openIDConnect"`
	OAuth         OAuth          `mapstructure:"oauth"`
	Boards        Boards         `mastructure:"boards"`
	CustomBoards  []CustomBoard  `mapstructure:"customBoards"`
	Local         bool           `mapstructure:"local"`
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
