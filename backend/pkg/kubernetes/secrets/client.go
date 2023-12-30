package secrets

import (
	"context"
	"strconv"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Values struct {
	Host        string `json:"api" mapstructure:"api"`
	KyvernoAPI  string `json:"kyvernoApi" mapstructure:"kyvernoApi"`
	Certificate string `json:"certificate" mapstructure:"certificate"`
	SkipTLS     bool   `json:"skipTLS" mapstructure:"skipTLS"`
	Username    string `json:"username" mapstructure:"username"`
	Password    string `json:"password" mapstructure:"password"`

	Domain       string `json:"domain" mapstructure:"domain"`
	ClientID     string `json:"clientId" mapstructure:"clientId"`
	ClientSecret string `json:"clientSecret" mapstructure:"clientSecret"`
}

type Client interface {
	Get(context.Context, string) (Values, error)
}

type k8sClient struct {
	client v1.SecretInterface
}

func (c *k8sClient) Get(ctx context.Context, name string) (Values, error) {
	secret, err := kubernetes.Retry(func() (*corev1.Secret, error) {
		return c.client.Get(ctx, name, metav1.GetOptions{})
	})

	values := Values{}
	if err != nil {
		return values, err
	}

	if api, ok := secret.Data["api"]; ok {
		values.Host = string(api)
	}

	if kyvernoAPI, ok := secret.Data["kyvernoApi"]; ok {
		values.KyvernoAPI = string(kyvernoAPI)
	}

	if certificate, ok := secret.Data["certificate"]; ok {
		values.Certificate = string(certificate)
	}

	if username, ok := secret.Data["username"]; ok {
		values.Username = string(username)
	}

	if password, ok := secret.Data["password"]; ok {
		values.Password = string(password)
	}

	if domain, ok := secret.Data["domain"]; ok {
		values.Domain = string(domain)
	}

	if clientID, ok := secret.Data["clientId"]; ok {
		values.ClientID = string(clientID)
	}

	if clientSecret, ok := secret.Data["clientSecret"]; ok {
		values.ClientSecret = string(clientSecret)
	}

	if skipTLS, ok := secret.Data["skipTLS"]; ok {
		v, err := strconv.ParseBool(string(skipTLS))
		if err != nil {
			zap.L().Error("failed to parse skipTLS", zap.Error(err), zap.String("secret", name))
		} else {
			values.SkipTLS = v
		}
	}

	return values, nil
}

func NewClient(secretClient v1.SecretInterface) Client {
	return &k8sClient{secretClient}
}
