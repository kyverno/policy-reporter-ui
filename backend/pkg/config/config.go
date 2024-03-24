package config

import (
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"github.com/kyverno/policy-reporter-ui/pkg/logging"
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

// APISetup configuration
type Cluster struct {
	Name        string    `mapstructure:"name"`
	Host        string    `mapstructure:"host"`
	Plugins     []Plugin  `mapstructure:"plugins"`
	SkipTLS     bool      `mapstructure:"skipTLS"`
	Certificate string    `mapstructure:"certificate"`
	SecretRef   string    `mapstructure:"secretRef"`
	BasicAuth   BasicAuth `mapstructure:"basicAuth"`
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

type DefaultFilter struct {
	Resources        []string `mapstructure:"resources"`
	ClusterResources []string `mapstructure:"clusterResources"`
}

type UI struct {
	Path          string        `mapstructure:"path"`
	Disabled      bool          `mapstructure:"enabled"`
	DisplayMode   string        `mapstructure:"displayMode"`
	LabelFilter   []string      `mapstructure:"labelFilter"`
	DefaultFilter DefaultFilter `mapstructure:"defaultFilter"`
}

type Server struct {
	Port          int  `mapstructure:"port"`
	CORS          bool `mapstructure:"cors"`
	Debug         bool `mapstructure:"debug"`
	Logging       bool `mapstructure:"logging"`
	OverwriteHost bool `mapstructure:"overwriteHost"`
}

type Source struct {
	Name       string `mapstructure:"name"`
	Exceptions bool   `mapstructure:"exceptions"`
	Excludes   struct {
		NamespaceKinds []string `mapstructure:"namespaceKinds"`
		ClusterKinds   []string `mapstructure:"clusterKinds"`
		Results        []string `mapstructure:"results"`
	} `mapstructure:"excludes"`
}

type CustomBoard struct {
	Name       string `mapstructure:"name"`
	Namespaces struct {
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
	CustomBoards  []CustomBoard  `mapstructure:"customBoards"`
	Local         bool           `mapstructure:"local"`
}

func (c *Config) AuthEnabled() bool {
	return c.OAuth.Enabled || c.OpenIDConnect.Enabled
}
