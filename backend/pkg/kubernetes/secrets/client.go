package secrets

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

type Plugin struct {
	Name        string `json:"name" mapstructure:"name"`
	Host        string `json:"host" mapstructure:"host"`
	Certificate string `json:"certificate" mapstructure:"certificate"`
	SkipTLS     bool   `json:"skipTLS" mapstructure:"skipTLS"`
	Username    string `json:"username" mapstructure:"username"`
	Password    string `json:"password" mapstructure:"password"`
}

type Values struct {
	Host        string   `json:"host" mapstructure:"host"`
	Plugins     []Plugin `json:"plugins" mapstructure:"plugins"`
	Certificate string   `json:"certificate" mapstructure:"certificate"`
	SkipTLS     bool     `json:"skipTLS" mapstructure:"skipTLS"`
	Username    string   `json:"username" mapstructure:"username"`
	Password    string   `json:"password" mapstructure:"password"`
	SecretRef   string   `json:"secretRef" mapstructure:"secretRef"`

	// OAuth Values
	Provider string `json:"provider" mapstructure:"provider"`
	// OpenIDConnect
	DiscoveryURL string `json:"domain" mapstructure:"discoveryUrl"`
	// OAuth + OpenIDConnect
	ClientID     string `json:"clientId" mapstructure:"clientId"`
	ClientSecret string `json:"clientSecret" mapstructure:"clientSecret"`
}

func (v Values) Merge(n Values) Values {
	v.Username = utils.Fallback(v.Username, n.Username)
	v.Password = utils.Fallback(v.Password, n.Password)

	plugins := make([]Plugin, 0, len(v.Plugins))
	for _, p := range v.Plugins {
		p.Username = utils.Fallback(p.Username, n.Username)
		p.Password = utils.Fallback(p.Password, n.Password)

		plugins = append(plugins, p)
	}

	v.Plugins = plugins

	return v
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

	values := Values{
		Plugins: make([]Plugin, 0),
	}

	if err != nil {
		return values, err
	}

	if host, ok := secret.Data["host"]; ok {
		values.Host = string(host)
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

	if secretRef, ok := secret.Data["secretRef"]; ok {
		values.SecretRef = string(secretRef)
	}

	if domain, ok := secret.Data["domain"]; ok {
		values.DiscoveryURL = string(domain)
	}

	if domain, ok := secret.Data["discoveryUrl"]; ok {
		values.DiscoveryURL = string(domain)
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

	for k, v := range secret.Data {
		if !strings.HasPrefix(k, "plugin.") {
			continue
		}

		plugin := Plugin{}
		if err := json.Unmarshal(v, &plugin); err != nil {
			zap.L().Error("failed to unmarshal plugin config", zap.Error(err), zap.String("plugin", k))
			continue
		}

		values.Plugins = append(values.Plugins, plugin)
	}

	return values, nil
}

func NewClient(secretClient v1.SecretInterface) Client {
	return &k8sClient{secretClient}
}
