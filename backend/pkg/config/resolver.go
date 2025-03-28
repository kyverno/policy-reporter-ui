package config

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"go.uber.org/zap"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/kyverno/policy-reporter-ui/pkg/api"
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/api/proxy"
	"github.com/kyverno/policy-reporter-ui/pkg/auth"
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
	config    *Config
	secrets   secrets.Client
	k8sConfig *rest.Config
	clientset *k8s.Clientset
}

func (r *Resolver) CoreClient(cluster Cluster) (*core.Client, error) {
	options := []api.ClientOption{
		api.WithBaseURL(cluster.Host),
	}

	if cluster.Certificate != "" {
		options = append(options, api.WithCertificate(cluster.Certificate))
	} else if cluster.SkipTLS {
		options = append(options, api.WithSkipTLS())
	}

	if cluster.BasicAuth.Username != "" {
		options = append(options, api.WithBaseAuth(api.BasicAuth{
			Username: cluster.BasicAuth.Username,
			Password: cluster.BasicAuth.Password,
		}))
	}

	if r.config.Logging.API {
		options = append(options, api.WithLogging())
	}

	return core.New(options)
}

func (r *Resolver) PluginClient(p Plugin) (*plugin.Client, error) {
	options := []api.ClientOption{
		api.WithBaseURL(p.Host),
	}

	if p.Certificate != "" {
		options = append(options, api.WithCertificate(p.Certificate))
	} else if p.SkipTLS {
		options = append(options, api.WithSkipTLS())
	}

	if p.BasicAuth.Username != "" {
		options = append(options, api.WithBaseAuth(api.BasicAuth{
			Username: p.BasicAuth.Username,
			Password: p.BasicAuth.Password,
		}))
	}

	if r.config.Logging.API {
		options = append(options, api.WithLogging())
	}

	return plugin.New(options)
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

func (r *Resolver) LoadPluginSecret(ctx context.Context, plugin Plugin) (Plugin, error) {
	if plugin.SecretRef != "" {
		values, err := r.LoadSecret(ctx, plugin.SecretRef)
		if err != nil {
			return plugin, err
		}

		plugin = plugin.FromValues(values)
	}

	return plugin, nil
}

func (r *Resolver) Proxy(cluster Cluster) (*httputil.ReverseProxy, error) {
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

	if r.config.Logging.API {
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

	return proxy.New(target, options, proxyOptions), nil
}

func (r *Resolver) LoadSecret(ctx context.Context, secretRef string) (secrets.Values, error) {
	if r.secrets == nil {
		clientset, err := r.Clientset()
		if err != nil {
			return secrets.Values{}, err
		}

		zap.L().Debug("create secret client", zap.String("namespace", r.config.Namespace))
		r.secrets = secrets.NewClient(clientset.CoreV1().Secrets(r.config.Namespace))
	}

	values, err := r.secrets.Get(ctx, secretRef)
	if err != nil {
		return secrets.Values{}, err
	}

	zap.L().Info("values loaded from secret", zap.String("secretRef", secretRef))

	if values.SecretRef != "" {
		nested, err := r.secrets.Get(ctx, values.SecretRef)
		if err != nil {
			return secrets.Values{}, err
		}
		values = values.Merge(nested)
		zap.L().Info("auth values loaded from secret", zap.String("secretRef", values.SecretRef))
	}

	return values, nil
}

func (r *Resolver) K8sConfig() (*rest.Config, error) {
	if r.k8sConfig != nil {
		return r.k8sConfig, nil
	}

	var k8sConfig *rest.Config
	var err error

	if r.config.Local {
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
	config := r.config.OAuth

	if config.SecretRef != "" {
		values, err := r.LoadSecret(ctx, config.SecretRef)
		if err != nil {
			return nil, err
		}

		config = config.FromValues(values)
	}

	provider := auth.NewProvider(config.Provider, config.ClientID, config.ClientSecret, config.CallbackURL, config.Scopes)
	if provider == nil {
		return nil, errors.New("provider not available")
	}

	goth.UseProviders(provider)
	auth.Setup(engine, r.config.OAuth.BasePath(), r.config.AuthGroupClaim(), config.Provider, r.config.TempDir)

	return []gin.HandlerFunc{auth.Provider(r.config.OAuth.Provider), auth.Auth(r.config.OAuth.BasePath())}, nil
}

func (r *Resolver) SetupOIDC(ctx context.Context, engine *gin.Engine) ([]gin.HandlerFunc, error) {
	oid := r.config.OpenIDConnect

	if oid.SecretRef != "" {
		values, err := r.LoadSecret(ctx, oid.SecretRef)
		if err != nil {
			zap.L().Error("failed to load openIDConnect secret", zap.String("secret", oid.SecretRef), zap.Error(err))
			return nil, err
		}

		oid = oid.FromValues(values)
	}

	zap.L().Debug("setup openIDConnect", zap.String("callback", oid.Callback()), zap.String("discovery", oid.Discovery()))

	client := api.NewHTTPClient()
	if oid.Certificate != "" {
		pool, err := api.LoadCerts(oid.Certificate)
		if err != nil {
			zap.L().Error("failed to load certificate for OIDC provider", zap.Error(err), zap.String("path", r.config.OpenIDConnect.Certificate))
			return nil, err
		}
		client.Transport.(*http.Transport).TLSClientConfig.RootCAs = pool
	}
	if oid.SkipTLS {
		client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = true
	}

	provider, err := openidConnect.New(oid.ClientID, oid.ClientSecret, oid.Callback(), oid.Discovery(), client, oid.Scopes...)
	if err != nil {
		zap.L().Error("failed to create openIDConnect provider", zap.Error(err))
		return nil, err
	}

	goth.UseProviders(provider)

	auth.Setup(engine, r.config.OpenIDConnect.BasePath(), r.config.AuthGroupClaim(), "openid-connect", r.config.TempDir)

	return []gin.HandlerFunc{auth.Provider("openid-connect"), auth.Auth(r.config.OpenIDConnect.BasePath())}, nil
}

func (r *Resolver) Server(ctx context.Context) (*server.Server, error) {
	if !r.config.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	if r.config.Server.CORS {
		engine.Use(cors.Default())
	}

	middleware := []gin.HandlerFunc{}

	if r.config.OpenIDConnect.Enabled {
		handler, err := r.SetupOIDC(ctx, engine)
		if err != nil {
			zap.L().Error("failed to setup oidc", zap.Error(err))
			return nil, err
		}

		middleware = append(middleware, handler...)
	} else if r.config.OAuth.Enabled {
		handler, err := r.SetupOAuth(ctx, engine)
		if err != nil {
			zap.L().Error("failed to setup oauth", zap.Error(err))
			return nil, err
		}

		middleware = append(middleware, handler...)
	}

	if r.config.Logging.Server {
		middleware = append(
			middleware,
			ginzap.Ginzap(r.Logger(), time.RFC3339, true),
			ginzap.RecoveryWithZap(r.Logger(), true),
		)
	} else {
		middleware = append(middleware, gin.Recovery())
	}

	if r.config.AuthEnabled() {
		middleware = append(middleware, auth.ClusterPermissions(MapClusterPermissions(r.config)))
	}

	serv := server.NewServer(engine, r.config.Server.Port, middleware)

	for _, cluster := range r.config.Clusters {
		cluster, err := r.LoadClusterSecret(ctx, cluster)
		if err != nil {
			zap.L().Error("failed to load cluster secret", zap.Error(err), zap.String("cluser", cluster.Name), zap.String("secretRef", cluster.SecretRef))
			continue
		}

		proxy, err := r.Proxy(cluster)
		if err != nil {
			zap.L().Error("failed to resolve proxies", zap.Error(err), zap.String("cluser", cluster.Name), zap.String("host", cluster.Host))
			continue
		}

		client, err := r.CoreClient(cluster)
		if err != nil {
			zap.L().Error("failed to create core api client", zap.Error(err), zap.String("cluser", cluster.Name), zap.String("host", cluster.Host))
			continue
		}

		plugins := make(map[string]*plugin.Client, len(cluster.Plugins))
		for _, p := range cluster.Plugins {
			p, err := r.LoadPluginSecret(ctx, p)
			if err != nil {
				zap.L().Error(
					"failed to load plugin secret",
					zap.Error(err),
					zap.String("cluster", cluster.Name),
					zap.String("plugin", p.Name),
					zap.String("secretRef", p.SecretRef),
				)
				continue
			}

			pClient, err := r.PluginClient(p)
			if err != nil {
				zap.L().Error("failed to create plugin client", zap.Error(err), zap.String("cluser", cluster.Name), zap.String("plugin", p.Name))
				continue
			}

			plugins[p.Name] = pClient
		}

		serv.RegisterCluster(cluster.Name, client, plugins, proxy)
	}

	if !r.config.UI.Disabled {
		var uiMiddleware []gin.HandlerFunc
		if r.config.AuthEnabled() {
			uiMiddleware = append(uiMiddleware, auth.Valid(r.config.AuthBasePath(), r.config.AuthGroupClaim()))
		}

		zap.L().Info("register UI", zap.String("path", r.config.UI.Path))
		serv.RegisterUI(r.config.UI.Path, uiMiddleware)
	}

	serv.RegisterAPI(MapConfig(r.config), MapCustomBoards(r.config.CustomBoards))

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
