package config

import (
	"context"
	"errors"
	"net/http/httputil"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"github.com/kyverno/policy-reporter-ui/pkg/logging"
	"github.com/kyverno/policy-reporter-ui/pkg/proxy"
	"github.com/kyverno/policy-reporter-ui/pkg/server"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
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

func (r *Resolver) ExternalProxies(ctx context.Context, cluster Cluster) (map[string]*httputil.ReverseProxy, error) {
	if cluster.SecretRef != "" {
		values, err := r.LoadSecret(ctx, cluster.SecretRef)
		if err != nil {
			return nil, err
		}

		cluster = cluster.FromValues(values)
	}

	if cluster.Host == "" {
		return nil, ErrMissingAPI
	}

	target, err := url.Parse(cluster.Host)
	if err != nil {
		return nil, err
	}

	options := make([]proxy.DirectorOption, 0)
	proxyOptions := make([]proxy.ProxyOption, 0)
	basicAuth := cluster.BasicAuth

	if r.config.Logging.Enabled {
		options = append(options, proxy.WithLogging())
	}

	if r.config.Server.OverwriteHost {
		options = append(options, proxy.WithHostOverwrite())
	}

	if basicAuth.Username != "" && basicAuth.Password != "" {
		options = append(options, proxy.WithAuth(basicAuth.Username, basicAuth.Password))
	}

	if cluster.SkipTLS {
		proxyOptions = append(proxyOptions, proxy.WithSkipTLS())
	}

	if cluster.Certificate != "" {
		proxyOptions = append(proxyOptions, proxy.WithCertificate(cluster.Certificate))
	}

	proxies := map[string]*httputil.ReverseProxy{
		"core": proxy.New(target, options, proxyOptions),
	}

	for _, p := range cluster.Plugins {
		pluginTarget, err := url.Parse(p.Host)
		if err != nil {
			zap.L().Error("failed to parse plugin host", zap.String("plugin", p.Name), zap.String("host", p.Host), zap.Error(err))
			continue
		}

		proxies[p.Name] = proxy.New(pluginTarget, options, proxyOptions)
	}

	return proxies, nil
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

func (r *Resolver) InitSecretClient() error {
	var k8sConfig *rest.Config
	var err error

	if r.config.KubeConfig.CurrentContext != "" {
		k8sConfig, err = utils.RestConfig(r.config.KubeConfig)
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

func (r *Resolver) Server(ctx context.Context) *server.Server {
	gin.SetMode(r.config.Server.Mode)

	engine := gin.Default()
	if r.config.Server.CORS {
		engine.Use(cors.Default())
	}

	serv := server.NewServer(engine, r.config.Server.Port)

	for _, cluster := range r.config.Clusters {
		proxies, err := r.ExternalProxies(ctx, cluster)
		if err != nil {
			zap.L().Error("failed to resolve proxies", zap.Error(err), zap.String("cluser", cluster.Name))
			continue
		}

		serv.RegisterCluster(cluster.Name, proxies)
	}

	if !r.config.UI.Disabled {
		zap.L().Info("register UI", zap.String("path", r.config.UI.Path))
		serv.RegisterUI(r.config.UI.Path)
	}

	serv.RegisterAPI(MapConfig(r.config))

	return serv
}

func (r *Resolver) Logger() *zap.Logger {
	return logging.New(r.config.Logging)
}

func NewResolver(config *Config) *Resolver {
	return &Resolver{
		config: config,
	}
}
