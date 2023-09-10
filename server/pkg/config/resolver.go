package config

import (
	"context"
	"errors"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"github.com/kyverno/policy-reporter-ui/pkg/logging"
	"github.com/kyverno/policy-reporter-ui/pkg/proxy"
	"github.com/kyverno/policy-reporter-ui/pkg/redis"
	"github.com/kyverno/policy-reporter-ui/pkg/report"
)

var (
	ErrMissingClient = errors.New("secret client was not initialized")
	ErrMissingAPI    = errors.New("missing core api configuration")
)

type Resolver struct {
	config     *Config
	kubeConfig string
	devMode    bool

	secrets secrets.Client
}

func (r *Resolver) InternalProxies(ctx context.Context, core, plugin string) (*httputil.ReverseProxy, *httputil.ReverseProxy, error) {
	target, err := url.Parse(core)
	if err != nil {
		return nil, nil, err
	}

	options := make([]proxy.DirectorOption, 0)

	apiConfig := r.config.APIConfig
	basicAuth := &apiConfig.BasicAuth

	if apiConfig.Logging {
		options = append(options, proxy.WithLogging())
		zap.L().Info("api logging enabled")
	}

	if apiConfig.OverwriteHost {
		options = append(options, proxy.WithHostOverwrite())
		zap.L().Info("host overwrite enabled")
	}

	if basicAuth.SecretRef != "" {
		basicAuth, err = r.LoadBasicAuth(ctx, basicAuth.SecretRef)
		if err != nil {
			return nil, nil, err
		}
		zap.L().Info("authentication loaded from secret", zap.String("secretRef", basicAuth.SecretRef))
	}

	if basicAuth.Username != "" && basicAuth.Password != "" {
		options = append(options, proxy.WithAuth(basicAuth.Username, basicAuth.Password))
	}

	coreProxy := proxy.New(target, options, make([]proxy.ProxyOption, 0))

	var pluginProxy *httputil.ReverseProxy
	if plugin != "" {
		pluginTarget, err := url.Parse(plugin)
		if err != nil {
			return nil, nil, err
		}

		pluginProxy = proxy.New(pluginTarget, options, make([]proxy.ProxyOption, 0))
	}

	return coreProxy, pluginProxy, nil
}

func (r *Resolver) ExternalProxies(ctx context.Context, api API) (*httputil.ReverseProxy, *httputil.ReverseProxy, error) {
	if api.SecretRef != "" {
		values, err := r.LoadSecret(ctx, api.SecretRef)
		if err != nil {
			return nil, nil, err
		}

		api = api.FromValues(values)
	}

	if api.CoreAPI == "" {
		return nil, nil, ErrMissingAPI
	}

	target, err := url.Parse(api.CoreAPI)
	if err != nil {
		return nil, nil, err
	}

	options := make([]proxy.DirectorOption, 0)
	proxyOptions := make([]proxy.ProxyOption, 0)
	basicAuth := api.BasicAuth

	if r.config.APIConfig.Logging {
		options = append(options, proxy.WithLogging())
	}

	if r.config.APIConfig.OverwriteHost {
		options = append(options, proxy.WithHostOverwrite())
	}

	if basicAuth.Username != "" && basicAuth.Password != "" {
		options = append(options, proxy.WithAuth(basicAuth.Username, basicAuth.Password))
	}

	if api.SkipTLS {
		proxyOptions = append(proxyOptions, proxy.WithSkipTLS())
	}

	if api.Certificate != "" {
		proxyOptions = append(proxyOptions, proxy.WithCertificate(api.Certificate))
	}

	coreProxy := proxy.New(target, options, proxyOptions)

	var pluginProxy *httputil.ReverseProxy
	if api.KyvernoAPI != "" {
		pluginTarget, err := url.Parse(api.KyvernoAPI)
		if err != nil {
			return nil, nil, err
		}

		pluginProxy = proxy.New(pluginTarget, options, proxyOptions)
	}

	return coreProxy, pluginProxy, nil
}

func (r *Resolver) LoadBasicAuth(ctx context.Context, secretRef string) (*BasicAuth, error) {
	if r.secrets == nil {
		return nil, ErrMissingClient
	}

	values, err := r.secrets.Get(ctx, secretRef)
	if err != nil {
		return nil, err
	}

	return &BasicAuth{
		Username:  values.Username,
		Password:  values.Password,
		SecretRef: secretRef,
	}, nil
}

func (r *Resolver) LoadSecret(ctx context.Context, secretRef string) (secrets.Values, error) {
	if r.secrets == nil {
		return secrets.Values{}, ErrMissingClient
	}

	values, err := r.secrets.Get(ctx, secretRef)
	if err != nil {
		return secrets.Values{}, err
	}

	zap.L().Info("values loaded from secret", zap.String("secretRef", secretRef))
	return values, nil
}

func (r *Resolver) InitClusters(enablePlugin bool) {
	r.config.Clusters = make([]Cluster, 0, len(r.config.APIs)+1)
	if len(r.config.APIs) > 0 {
		if r.config.ClusterName == "" {
			r.config.ClusterName = "Default"
		}

		r.config.Clusters = append(r.config.Clusters, Cluster{
			Name:    r.config.ClusterName,
			Kyverno: enablePlugin,
		})
	}
}

func (r *Resolver) InitSecretClient() error {
	var k8sConfig *rest.Config
	var err error

	if r.kubeConfig != "" {
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", r.kubeConfig)
	} else {
		k8sConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return err
	}

	clientset, err := k8s.NewForConfig(k8sConfig)
	if err != nil {
		return err
	}

	zap.L().Info("secret client initialized")
	r.secrets = secrets.NewClient(clientset.CoreV1().Secrets(r.config.Namespace))

	return nil
}

func (r *Resolver) Store() report.Store {
	if r.config.Redis.Enabled {
		zap.L().Info("redis store enabled")
		return redis.NewFromConfig(r.config.Redis)
	}

	zap.L().Info("log store size", zap.Int("size", r.config.LogSize))
	return report.NewResultStore(r.config.LogSize)
}

func (r *Resolver) Logger() *zap.Logger {
	return logging.New(r.config.Logging)
}

func NewResolver(config *Config, kubeConfig string, devMode bool) *Resolver {
	return &Resolver{
		config:     config,
		kubeConfig: kubeConfig,
		devMode:    devMode,
	}
}
