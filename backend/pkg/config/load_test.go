package config

import (
	"os"
	"testing"
)

func TestLoadConfigFileErrors(t *testing.T) {
	for _, explicit := range []bool{false, true} {
		for _, state := range []string{"missing", "malformed", "unreadable", "valid"} {
			name := state
			if explicit {
				name = "explicit/" + name
			} else {
				name = "default/" + name
			}
			t.Run(name, func(t *testing.T) {
				t.Chdir(t.TempDir())
				filename := "config.yaml"
				switch state {
				case "malformed":
					if err := os.WriteFile(filename, []byte("server: ["), 0o600); err != nil {
						t.Fatal(err)
					}
				case "unreadable":
					// A directory reliably fails reading even when tests run as root.
					if err := os.Mkdir(filename, 0o700); err != nil {
						t.Fatal(err)
					}
				case "valid":
					if err := os.WriteFile(filename, []byte("server:\n  tls:\n    enabled: true\n"), 0o600); err != nil {
						t.Fatal(err)
					}
				}
				cfgFile := ""
				if explicit {
					cfgFile = filename
				}
				var config Config
				err := Load(&config, cfgFile)
				wantError := state != "valid" && (explicit || state != "missing")
				if (err != nil) != wantError {
					t.Fatalf("Load() error = %v, want error = %v", err, wantError)
				}
				if state == "valid" && !config.Server.TLS.Enabled {
					t.Fatal("TLS configuration was not loaded")
				}
			})
		}
	}
}
