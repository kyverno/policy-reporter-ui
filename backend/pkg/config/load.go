package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func Load(c *Config, cfgFile string) error {
	var k = koanf.New("!")

	if cfgFile == "" {
		cfgFile = "./config.yaml"
	}

	if err := k.Load(file.Provider(cfgFile), yaml.Parser()); err != nil {
		fmt.Printf("[ERROR] failed to load config file: %v\n", err)
	}

	k.Load(env.Provider("PR_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "PR_")), "_", ".", -1)
	}), nil)

	err := k.Unmarshal("", c)

	if c.UI.Path == "" {
		c.UI.Path = path.Join(os.Getenv("KO_DATA_PATH"), "ui")
	}

	return err
}
