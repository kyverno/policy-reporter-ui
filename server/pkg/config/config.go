package config

import (
	"log"

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

// API configuration
type API struct {
	Name       string `json:"name" mapstructure:"name"`
	CoreAPI    string `json:"api" mapstructure:"api"`
	KyvernoAPI string `json:"kyvernoApi" mapstructure:"kyvernoApi"`
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

// Config structure
type Config struct {
	Views       Views     `json:"views" mapstructure:"views"`
	LogSize     int       `json:"-" mapstructure:"logSize"`
	DisplayMode string    `json:"displayMode" mapstructure:"displayMode"`
	Plugins     []string  `json:"plugins" mapstructure:"-"`
	Clusters    []Cluster `json:"clusters" mapstructure:"-"`
	APIs        []API     `json:"-" mapstructure:"clusters"`
	Redis       Redis     `json:"-" mapstructure:"redis"`
}

// LoadConfig from config file
func LoadConfig(cfgFile string) (*Config, error) {
	v := viper.New()
	v.SetDefault("logSize", 500)

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

	return c, err
}
