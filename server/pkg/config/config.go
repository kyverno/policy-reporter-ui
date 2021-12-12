package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	LogSize     int      `json:"-" mapstructure:"logSize"`
	DisplayMode string   `json:"displayMode" mapstructure:"displayMode"`
	Plugins     []string `json:"plugins" mapstructure:"-"`
}

func NewConfig() *Config {
	return &Config{
		Plugins:     make([]string, 0),
		DisplayMode: "",
	}
}

func LoadConfig(cfgFile string) (*Config, error) {
	v := viper.New()
	v.SetDefault("logSize", 500)

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
