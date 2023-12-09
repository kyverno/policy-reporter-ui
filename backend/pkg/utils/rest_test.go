package utils_test

import (
	"testing"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

func TestRestConfig(t *testing.T) {
	tests := []struct {
		name       string
		kubeConfig string
		overrides  clientcmd.ConfigOverrides
		wantErr    bool
	}{{
		name:       "empty",
		kubeConfig: "../../testdata/.kube/config",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("KUBECONFIG", tt.kubeConfig)
			_, err := utils.RestConfig(tt.overrides)
			if (err != nil) != tt.wantErr {
				t.Errorf("RestConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
