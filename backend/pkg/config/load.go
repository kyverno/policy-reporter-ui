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
	k := koanf.New("!")

	if cfgFile == "" {
		cfgFile = "./config.yaml"
	}

	if err := k.Load(file.Provider(cfgFile), yaml.Parser()); err != nil {
		fmt.Printf("[ERROR] failed to load config file: %v\n", err)
	}

	err := k.Load(env.Provider("PR_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "PR_")), "_", ".", -1)
	}), nil)
	if err != nil {
		return err
	}

	err = k.Unmarshal("", c)
	if err != nil {
		return err
	}

	if c.UI.Path == "" {
		c.UI.Path = path.Join(os.Getenv("KO_DATA_PATH"), "ui")
	}

	if c.Namespace == "" {
		c.Namespace = os.Getenv("POD_NAMESPACE")
	}

	return err
}
