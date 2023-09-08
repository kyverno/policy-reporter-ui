package config

import (
	"log"
	"os"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"github.com/spf13/viper"
)

// Dashboard configuration
type Dashboard struct {
	PolicyReports        bool `json:"policyReports" mapstructure:"policyReports"`
	ClusterPolicyReports bool `json:"clusterPolicyReports" mapstructure:"clusterPolicyReports"`
}

// Views configuration
type Views struct {
	Dashboard            Dashboard `json:"dashboard" mapstructure:"dashboard"`
	Logs                 bool      `json:"logs" mapstructure:"logs"`
	PolicyReports        bool      `json:"policyReports" mapstructure:"policyReports"`
	ClusterPolicyReports bool      `json:"clusterPolicyReports" mapstructure:"clusterPolicyReports"`
	KyvernoPolicies      bool      `json:"kyvernoPolicies" mapstructure:"kyvernoPolicies"`
	KyvernoVerifyImages  bool      `json:"kyvernoVerifyImages" mapstructure:"kyvernoVerifyImages"`
}

// BasicAuth configuration
type BasicAuth struct {
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	SecretRef string `mapstructure:"secretRef"`
}

// API configuration
type API struct {
	Name        string    `json:"name" mapstructure:"name"`
	CoreAPI     string    `json:"api" mapstructure:"api"`
	KyvernoAPI  string    `json:"kyvernoApi" mapstructure:"kyvernoApi"`
	SkipTLS     bool      `json:"-" mapstructure:"skipTLS"`
	Certificate string    `json:"-" mapstructure:"certificate"`
	SecretRef   string    `json:"-" mapstructure:"secretRef"`
	BasicAuth   BasicAuth `json:"-" mapstructure:"basicAuth"`
}

func (a API) FromValues(values secrets.Values) API {
	if values.CoreAPI != "" {
		a.CoreAPI = values.CoreAPI
	}
	if values.KyvernoAPI != "" {
		a.KyvernoAPI = values.KyvernoAPI
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
	Name    string `json:"name" mapstructure:"name"`
	ID      string `json:"id" mapstructure:"id"`
	Kyverno bool   `json:"kyverno" mapstructure:"kyverno"`
}

// Redis configuration
type Redis struct {
	Enabled  bool   `mapstructure:"enabled"`
	Address  string `mapstructure:"address"`
	Prefix   string `mapstructure:"prefix"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type Logging struct {
	LogLevel    int8   `mapstructure:"logLevel"`
	Encoding    string `mapstructure:"encoding"`
	Development bool   `mapstructure:"development"`
}

type APIConfig struct {
	Logging   bool      `mapstructure:"logging"`
	BasicAuth BasicAuth `mapstructure:"basicAuth"`
}

// Config structure
type Config struct {
	Views           Views     `json:"views" mapstructure:"views"`
	LogSize         int       `json:"-" mapstructure:"logSize"`
	DisplayMode     string    `json:"displayMode" mapstructure:"displayMode"`
	RefreshInterval int       `json:"refreshInterval" mapstructure:"refreshInterval"`
	Plugins         []string  `json:"plugins" mapstructure:"-"`
	Clusters        []Cluster `json:"clusters" mapstructure:"-"`
	ClusterName     string    `json:"-" mapstructure:"clusterName"`
	APIs            []API     `json:"-" mapstructure:"clusters"`
	Redis           Redis     `json:"-" mapstructure:"redis"`
	LabelFilter     []string  `json:"labelFilter" mapstructure:"labelFilter"`
	Logging         Logging   `json:"-" mapstructure:"logging"`
	APIConfig       APIConfig `json:"-" mapstructure:"apiConfig"`
	Namespace       string    `json:"namespace" mapstructure:"namespace"`
}

// LoadConfig from config file
func LoadConfig(cfgFile string) (*Config, error) {
	v := viper.New()

	v.SetDefault("logSize", 500)
	v.SetDefault("refreshInterval", 10000)
	v.SetDefault("clusterName", "Default")

	v.SetDefault("views.logs", true)
	v.SetDefault("views.policyReports", true)
	v.SetDefault("views.clusterPolicyReports", true)
	v.SetDefault("views.kyvernoPolicies", true)
	v.SetDefault("views.kyvernoVerifyImages", true)

	v.SetDefault("views.dashboard.policyReports", true)
	v.SetDefault("views.dashboard.clusterPolicyReports", true)

	c := &Config{}

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.AddConfigPath(".")
		v.SetConfigName("config")
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Println("[INFO] No configuration file found")
	}

	err := v.Unmarshal(c)

	if ns := os.Getenv("POD_NAMESPACE"); ns != "" {
		c.Namespace = ns
	}

	return c, err
}
