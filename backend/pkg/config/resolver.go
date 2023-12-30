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

	"github.com/kyverno/policy-reporter-ui/pkg/auth"
	"github.com/kyverno/policy-reporter-ui/pkg/core/client"
	"github.com/kyverno/policy-reporter-ui/pkg/core/proxy"
	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"github.com/kyverno/policy-reporter-ui/pkg/logging"
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

	k8sConfig *rest.Config
	clientset *k8s.Clientset
}

func (r *Resolver) CoreClient(cluster Cluster) (*client.Client, error) {
	options := []client.ClientOption{
		client.WithBaseURL(cluster.Host),
	}

	if cluster.Certificate != "" {
		options = append(options, client.WithCertificate(cluster.Certificate))
	} else if cluster.SkipTLS {
		options = append(options, client.WithSkipTLS())
	}

	if cluster.BasicAuth.Username != "" {
		options = append(options, client.WithBaseAuth(client.BasicAuth{
			Username: cluster.BasicAuth.Username,
			Password: cluster.BasicAuth.Password,
		}))
	}

	if r.config.Logging.Enabled {
		options = append(options, client.WithLogging())
	}

	return client.New(options)
}

func (r *Resolver) LoadClusterSecret(ctx context.Context, cluster Cluster) (Cluster, error) {
	if cluster.SecretRef != "" {
		values, err := r.LoadSecret(ctx, cluster.SecretRef)
		if err != nil {
			return cluster, err
		}

		cluster = cluster.FromValues(values)
	}

	return cluster, nil
}

func (r *Resolver) ExternalProxies(cluster Cluster) (map[string]*httputil.ReverseProxy, error) {
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
		clientset, err := r.Clientset()
		if err != nil {
			return secrets.Values{}, err
		}

		r.secrets = secrets.NewClient(clientset.CoreV1().Secrets(r.config.Namespace))
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

func (r *Resolver) SetupOAuth(ctx context.Context, engine *gin.Engine) ([]gin.HandlerFunc, error) {
	if !r.config.OAuth.Enabled {
		return make([]gin.HandlerFunc, 0), nil
	}

	oauth := r.config.OAuth

	if oauth.SecretRef != "" {
		values, err := r.LoadSecret(ctx, oauth.SecretRef)
		if err != nil {
			return nil, err
		}

		oauth = oauth.FromValues(values)
	}

	authenticator, err := auth.New(
		oauth.Domain,
		oauth.ClientID,
		oauth.ClientSecret,
		oauth.Redirect,
		oauth.Scopes,
	)

	if err != nil {
		return nil, err
	}

	auth.Setup(engine, authenticator)

	return []gin.HandlerFunc{auth.Auth}, nil
}

func (r *Resolver) Server(ctx context.Context) (*server.Server, error) {
	if !r.config.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	if r.config.Server.CORS {
		engine.Use(cors.Default())
	}

	middleware, err := r.SetupOAuth(ctx, engine)
	if err != nil {
		zap.L().Error("failed to setup oauth", zap.Error(err))
	}

	serv := server.NewServer(engine, r.config.Server.Port, middleware)

	for _, cluster := range r.config.Clusters {
		cluster, err := r.LoadClusterSecret(ctx, cluster)
		if err != nil {
			zap.L().Error("failed to load cluster secret", zap.Error(err), zap.String("cluser", cluster.Name))
			continue
		}

		proxies, err := r.ExternalProxies(cluster)
		if err != nil {
			zap.L().Error("failed to resolve proxies", zap.Error(err), zap.String("cluser", cluster.Name))
			continue
		}

		client, err := r.CoreClient(cluster)

		serv.RegisterCluster(cluster.Name, client, proxies)
	}

	if !r.config.UI.Disabled {
		zap.L().Info("register UI", zap.String("path", r.config.UI.Path))
		serv.RegisterUI(r.config.UI.Path)
	}

	serv.RegisterAPI(
		MapConfig(r.config),
		MapCustomBoards(r.config.CustomBoards),
	)

	return serv, nil
}

func (r *Resolver) Logger() *zap.Logger {
	return logging.New(r.config.Logging)
}

func NewResolver(config *Config) *Resolver {
	return &Resolver{
		config: config,
	}
}
