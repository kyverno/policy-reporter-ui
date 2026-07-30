package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

func TestRestConfig(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config")
	config := "apiVersion: v1\nclusters:\n- cluster:\n    server: https://example.invalid\n  name: test\ncontexts:\n- context:\n    cluster: test\n    user: test\n  name: test\ncurrent-context: test\nusers:\n- name: test\n  user: {}\n"
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("KUBECONFIG", configPath)
	if _, err := utils.RestConfig(clientcmd.ConfigOverrides{}); err != nil {
		t.Fatalf("RestConfig() error = %v, want nil", err)
	}
}
