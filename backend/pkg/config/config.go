package config

import (
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"github.com/kyverno/policy-reporter-ui/pkg/logging"
	"github.com/kyverno/policy-reporter-ui/pkg/redis"
)

// BasicAuth configuration
type BasicAuth struct {
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	SecretRef string `mapstructure:"secretRef"`
}

type Plugin struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
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
	if values.KyvernoAPI != "" {
		a.Plugins = append(a.Plugins, Plugin{Name: "kyverno", Host: values.KyvernoAPI})
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
	Port          int    `mapstructure:"port"`
	CORS          bool   `mapstructure:"cors"`
	Mode          string `mapstructure:"mode"`
	OverwriteHost bool   `mapstructure:"overwriteHost"`
}

type Source struct {
	Name     string `mapstructure:"name"`
	Excludes struct {
		NamespaceKinds []string `mapstructure:"namespaceKinds"`
		ClusterKinds   []string `mapstructure:"clusterKinds"`
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
}

// Config structure
type Config struct {
	KubeConfig   clientcmd.ConfigOverrides
	Namespace    string         `mapstructure:"namespace"`
	Clusters     []Cluster      `mapstructure:"clusters"`
	Sources      []Source       `mapstructure:"sources"`
	Server       Server         `mapstructure:"server"`
	UI           UI             `mapstructure:"ui"`
	Logging      logging.Config `mapstructure:"logging"`
	Redis        redis.Config   `mapstructure:"redis"`
	CustomBoards []CustomBoard  `mapstructure:"customBoards"`
	Cluster      bool           `mapstructure:"cluster"`
}
