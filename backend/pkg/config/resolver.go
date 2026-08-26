package config

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/zapr"
	"github.com/gosimple/slug"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	ctrserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/kyverno/policy-reporter-ui/pkg/api"
	"github.com/kyverno/policy-reporter-ui/pkg/auth"
	"github.com/kyverno/policy-reporter-ui/pkg/cluster"
	uiv1alpha1 "github.com/kyverno/policy-reporter-ui/pkg/crd/api/ui/v1alpha1"
	"github.com/kyverno/policy-reporter-ui/pkg/customboard"
	kcl "github.com/kyverno/policy-reporter-ui/pkg/kubernetes/cluster"
	kcb "github.com/kyverno/policy-reporter-ui/pkg/kubernetes/customboard"
	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"github.com/kyverno/policy-reporter-ui/pkg/server"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

var (
	ErrMissingClient = errors.New("secret client was not initialized")
	ErrMissingAPI    = errors.New("missing core api configuration")
)

type Resolver struct {
	config       *Config
	secrets      secrets.Client
	k8sConfig    *rest.Config
	clientset    *k8s.Clientset
	customBoards *customboard.Collection
	clusters     *cluster.Collection
	mgr          manager.Manager
}

func (r *Resolver) SecretClient() (secrets.Client, error) {
	if r.secrets == nil {
		clientset, err := r.Clientset()
		if err != nil {
			return nil, err
		}

		zap.L().Debug("create secret client", zap.String("namespace", r.config.Namespace))
		r.secrets = secrets.NewClient(clientset.CoreV1().Secrets(r.config.Namespace))
	}

	return r.secrets, nil
}

func (r *Resolver) LoadSecret(ctx context.Context, secretRef string) (secrets.Values, error) {
	client, err := r.SecretClient()
	if err != nil {
		return secrets.Values{}, err
	}

	return client.Get(ctx, secretRef)
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
		zap.L().Error("failed to create k8s config", zap.Error(err))
		return nil, err
	}

	clientset, err := k8s.NewForConfig(k8sConfig)
	if err != nil {
		zap.L().Error("failed to create k8s clientset", zap.Error(err))
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
	auth.Setup(engine, r.config.OAuth.BasePath(), r.config.AuthGroupClaim(), config.Provider, auth.SessionStorage{
		Storage: r.config.Server.Sessions.Storage,
		TempDir: r.config.Server.Sessions.TempDir,
	})

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

	client := auth.NewHTTPClient()
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

	client.Transport = api.NewLoggingRoundTripper(client.Transport)

	provider, err := openidConnect.New(oid.ClientID, oid.ClientSecret, oid.Callback(), oid.Discovery(), client, oid.PKCE, oid.Scopes...)
	if err != nil {
		zap.L().Error("failed to create openIDConnect provider", zap.Error(err))
		return nil, err
	}

	goth.UseProviders(provider)

	auth.Setup(engine, r.config.OpenIDConnect.BasePath(), r.config.AuthGroupClaim(), "openid-connect", auth.SessionStorage{
		Storage:  r.config.Server.Sessions.Storage,
		TempDir:  r.config.Server.Sessions.TempDir,
		Addr:     r.config.Server.Sessions.Addr,
		Database: r.config.Server.Sessions.Database,
		Username: r.config.Server.Sessions.Username,
		Password: r.config.Server.Sessions.Password,
	})

	return []gin.HandlerFunc{auth.Provider("openid-connect"), auth.Auth(r.config.OpenIDConnect.BasePath())}, nil
}

func (r *Resolver) ClusterLoader() *cluster.SecretLoader {
	secrets, err := r.SecretClient()
	if err != nil {
		return nil
	}

	return cluster.NewSecretLoader(secrets, r.config.Logging.API)
}

func (r *Resolver) ClusterCollection(ctx context.Context) *cluster.Collection {
	if r.clusters != nil {
		return r.clusters
	}

	loader := r.ClusterLoader()
	if loader == nil {
		r.clusters = cluster.NewCollection()
		return r.clusters
	}

	clusters := loader.LoadConfigs(ctx, r.config.Clusters)

	r.clusters = cluster.NewCollection(clusters...)
	return r.clusters
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
			ginzap.Ginzap(zap.L(), time.RFC3339, true),
			ginzap.RecoveryWithZap(zap.L(), true),
		)
	} else {
		middleware = append(middleware, gin.Recovery())
	}

	if r.config.AuthEnabled() {
		middleware = append(middleware, auth.ClusterPermissions(r.ClusterCollection(ctx)))
	}

	serv := server.NewServer(engine, r.config.Server.Port, middleware, r.ClusterCollection(ctx))

	if !r.config.UI.Disabled {
		var uiMiddleware []gin.HandlerFunc
		if r.config.AuthEnabled() {
			uiMiddleware = append(uiMiddleware, auth.Valid(r.config.AuthBasePath()))
		}

		zap.L().Info("register UI", zap.String("path", r.config.UI.Path))
		serv.RegisterUI(r.config.UI.Path, uiMiddleware)
	}

	serv.RegisterAPI(MapConfig(r.config), r.CustomBoards())

	return serv, nil
}

func (r *Resolver) Mgr() (manager.Manager, error) {
	if r.mgr != nil {
		return r.mgr, nil
	}

	k8sConfig, err := r.K8sConfig()
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	err = uiv1alpha1.Install(scheme)
	if err != nil {
		return nil, err
	}

	log.SetLogger(zapr.NewLogger(zap.L()))

	port := "0"
	if r.config.Server.Metrics.Enabled {
		port = fmt.Sprintf(":%d", r.config.Server.Metrics.Port)
	}

	mgr, err := ctrl.NewManager(k8sConfig, ctrl.Options{
		Scheme: scheme,
		Metrics: ctrserver.Options{
			BindAddress: port,
		},
	})
	if err != nil {
		return nil, err
	}
	r.mgr = mgr

	return r.mgr, nil
}

func (r *Resolver) CustomBoardController() (*kcb.Client, error) {
	mgr, err := r.Mgr()
	if err != nil {
		return nil, err
	}

	return kcb.NewClient(mgr, r.CustomBoards())
}

func (r *Resolver) ClusterController(ctx context.Context) (*kcl.Client, error) {
	mgr, err := r.Mgr()
	if err != nil {
		return nil, err
	}

	return kcl.NewClient(mgr, r.ClusterCollection(ctx), r.ClusterLoader())
}

func (r *Resolver) CustomBoards() *customboard.Collection {
	if r.customBoards != nil {
		return r.customBoards
	}

	r.customBoards = customboard.NewCollection(utils.Map(r.config.CustomBoards, func(c customboard.CustomBoard) *customboard.CustomBoard {
		c.ID = slug.Make(c.Name)
		c.Filter = customboard.MapFilterFields(c.Filter)

		return &c
	})...)

	return r.customBoards
}

func NewResolver(config *Config) *Resolver {
	return &Resolver{
		config: config,
	}
}
