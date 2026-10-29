package config

import (
	"errors"
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

	explicitConfig := cfgFile != ""
	if cfgFile == "" {
		cfgFile = "./config.yaml"
	}

	if err := k.Load(file.Provider(cfgFile), yaml.Parser()); err != nil {
		if explicitConfig || !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to load config file %s: %w", cfgFile, err)
		}
	}

	err := k.Load(env.Provider("PR_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "PR_")), "_", ".")
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
