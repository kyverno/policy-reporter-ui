package secrets_test

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
)

const secretName = "secret-values"

func newFakeClient() v1.SecretInterface {
	return fake.NewSimpleClientset(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "default",
		},
		Data: map[string][]byte{
			"api":         []byte("http://localhost:9200"),
			"kyvernoApi":  []byte("http://localhost:9200/kyverno"),
			"username":    []byte("username"),
			"password":    []byte("password"),
			"skipTLS":     []byte("true"),
			"certificate": []byte("certs"),
		},
	}).CoreV1().Secrets("default")
}

func Test_Client(t *testing.T) {
	client := secrets.NewClient(newFakeClient())

	t.Run("Get values from existing secret", func(t *testing.T) {
		values, err := client.Get(context.Background(), secretName)
		if err != nil {
			t.Errorf("Unexpected error while fetching secret: %s", err)
		}

		if values.CoreAPI != "http://localhost:9200" {
			t.Errorf("Unexpected CoreAPI: %s", values.CoreAPI)
		}

		if values.KyvernoAPI != "http://localhost:9200/kyverno" {
			t.Errorf("Unexpected KyvernoAPI: %s", values.KyvernoAPI)
		}

		if !values.SkipTLS {
			t.Errorf("Unexpected SkipTLS: %t", values.SkipTLS)
		}

		if values.Certificate != "certs" {
			t.Errorf("Unexpected Certificate: %s", values.Certificate)
		}

		if values.Username != "username" {
			t.Errorf("Unexpected Username: %s", values.Username)
		}

		if values.Password != "password" {
			t.Errorf("Unexpected Password: %s", values.Password)
		}
	})

	t.Run("Get values from not existing secret", func(t *testing.T) {
		_, err := client.Get(context.Background(), "not-exist")
		if !errors.IsNotFound(err) {
			t.Errorf("Expected not found error")
		}
	})
}
