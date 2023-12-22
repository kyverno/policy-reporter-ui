package config

import (
	"context"
	"errors"
	"net/http/httputil"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.uber.org/zap"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/namespaces"
	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"github.com/kyverno/policy-reporter-ui/pkg/logging"
	"github.com/kyverno/policy-reporter-ui/pkg/proxy"
	"github.com/kyverno/policy-reporter-ui/pkg/server"
	"github.com/kyverno/policy-reporter-ui/pkg/server/api"
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

	secrets    secrets.Client
	namespaces namespaces.Client

	k8sConfig *rest.Config
	clientset *k8s.Clientset
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

func (r *Resolver) K8sConfig() (*rest.Config, error) {
	if r.k8sConfig != nil {
		return r.k8sConfig, nil
	}

	var k8sConfig *rest.Config
	var err error

	if r.config.Cluster {
		k8sConfig, err = utils.RestConfig(r.config.KubeConfig)
	} else {
		k8sConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	r.k8sConfig = k8sConfig

	return r.k8sConfig, nil
}

func (r *Resolver) Clientset() (*k8s.Clientset, error) {
	if r.clientset != nil {
		return r.clientset, nil
	}

	k8sConfig, err := r.K8sConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := k8s.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	r.clientset = clientset

	return r.clientset, nil
}

func (r *Resolver) NamespaceClient() (namespaces.Client, error) {
	if r.namespaces != nil {
		return r.namespaces, nil
	}

	clientset, err := r.Clientset()
	if err != nil {
		return nil, err
	}

	r.namespaces = namespaces.NewClient(clientset.CoreV1().Namespaces())

	return r.namespaces, nil
}

func (r *Resolver) InitSecretClient() error {
	clientset, err := r.Clientset()
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

	if len(r.config.CustomBoards) > 0 {
		client, err := r.NamespaceClient()
		if err != nil {
			zap.L().Error("failed to create namespace client", zap.Error(err))

			return serv
		}

		configs := make(map[string]api.CustomBoard, len(r.config.CustomBoards))
		for _, c := range r.config.CustomBoards {
			id := slug.Make(c.Name)

			configs[id] = api.CustomBoard{
				Name: c.Name,
				ID:   id,
				Namespaces: api.Namespaces{
					Selector: c.Namespaces.Selector,
					List:     c.Namespaces.List,
				},
				Sources: api.Sources{
					List: c.Sources.List,
				},
				PolicyReports: api.PolicyReports{
					Selector: c.PolicyReports.Selector,
				}}
		}

		serv.RegisterCustomBoards(client, configs)
	}

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
