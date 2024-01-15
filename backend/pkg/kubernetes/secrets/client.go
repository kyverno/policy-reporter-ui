package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	k8s "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
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

	Domain       string `json:"domain" mapstructure:"domain"`
	ClientID     string `json:"clientId" mapstructure:"clientId"`
	ClientSecret string `json:"clientSecret" mapstructure:"clientSecret"`
}

const (
	Added = iota
	Updated
	Deleted
)

type Event struct {
	Values Values
	Type   int
}

type Client interface {
	Get(context.Context, string) (Values, error)
	Watch(ctx context.Context) (<-chan Event, error)
}

type k8sClient struct {
	client   v1.SecretInterface
	informer cache.SharedIndexInformer
	factory  informers.SharedInformerFactory
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

func mapping(secret *corev1.Secret) Values {
	values := Values{
		Plugins: make([]Plugin, 0),
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
			zap.L().Error("failed to parse skipTLS", zap.Error(err))
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

	return values
}

func (c *k8sClient) Watch(ctx context.Context) (<-chan Event, error) {
	channel := make(chan Event)
	c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if item, ok := obj.(*corev1.Secret); ok {
				values := mapping(item)
				channel <- Event{Values: values, Type: Added}

			}
		},
		UpdateFunc: func(_, obj interface{}) {
			if item, ok := obj.(*corev1.Secret); ok {
				values := mapping(item)
				channel <- Event{Values: values, Type: Updated}
			}
		},
		DeleteFunc: func(obj interface{}) {
			if item, ok := obj.(*corev1.Secret); ok {
				values := mapping(item)
				channel <- Event{Values: values, Type: Deleted}
			}
		},
	})

	c.factory.Start(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), c.informer.HasSynced) {
		return channel, fmt.Errorf("failed to sync secrets informer")
	}

	return channel, nil
}

func NewClient(ns string, clientset k8s.Interface) Client {
	client := clientset.CoreV1().Secrets(ns)

	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0,
		informers.WithNamespace(ns),
		informers.WithTweakListOptions(func(lo *metav1.ListOptions) {
			lo.LabelSelector = metav1.FormatLabelSelector(
				&metav1.LabelSelector{MatchLabels: map[string]string{
					"policy-reporter.config/type": "cluster",
				}},
			)
		}))

	informer := factory.Core().V1().Secrets().Informer()

	return &k8sClient{client, informer, factory}
}
