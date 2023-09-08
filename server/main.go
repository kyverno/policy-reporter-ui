package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"go.uber.org/zap"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kyverno/policy-reporter-ui/pkg/api"
	"github.com/kyverno/policy-reporter-ui/pkg/config"
	"github.com/kyverno/policy-reporter-ui/pkg/logging"
	"github.com/kyverno/policy-reporter-ui/pkg/proxy"
	"github.com/kyverno/policy-reporter-ui/pkg/redis"
	"github.com/kyverno/policy-reporter-ui/pkg/report"
)

var (
	kubeConfig     string
	configFile     string
	policyReporter string
	kyvernoPlugin  string
	overwriteHost  bool
	development    bool
	noUI           bool
	port           int
)

func main() {
	flag.StringVar(&configFile, "config", "./config.yaml", "Path to the config file")
	flag.StringVar(&policyReporter, "policy-reporter", "", "PolicyReporter Host")
	flag.StringVar(&kyvernoPlugin, "kyverno-plugin", "", "Kyverno Plugin Host")
	flag.StringVar(&kubeConfig, "k", "", "Kubernetes Config File")
	flag.BoolVar(&overwriteHost, "overwrite-host", false, "Overwrite Proxy Host and set Forward Header")
	flag.IntVar(&port, "port", 8080, "PolicyReporter UI port")
	flag.BoolVar(&development, "dev", false, "Enable CORS Header for development")
	flag.BoolVar(&noUI, "no-ui", false, "Disable the embedded frontend")
	flag.Parse()

	conf, err := config.LoadConfig(configFile)
	conf.Clusters = make([]config.Cluster, 0, len(conf.APIs)+1)
	if len(conf.APIs) > 0 {
		conf.Clusters = append(conf.Clusters, config.Cluster{
			Name:    conf.ClusterName,
			Kyverno: len(kyvernoPlugin) > 0,
		})
	}

	logger := logging.New(conf)

	if development {
		logger.Info("Development Mode enabled")
	}

	if overwriteHost {
		logger.Info("host overwrite enabled")
	}

	secretClient, err := secretClient(kubeConfig, conf.Namespace)
	if err != nil {
		logger.Warn("failed to setup secret client, secretRefs can not be resolved", zap.Error(err))
	}

	router := mux.NewRouter()

	var store report.Store

	if conf.Redis.Enabled {
		if development {
			logger.Info("redis store enabled")
		}

		store = redis.NewFromConfig(conf)
	} else {
		if development {
			logger.Info("log store size", zap.Int("size", conf.LogSize))
		}

		store = report.NewResultStore(conf.LogSize)
	}

	backend, err := url.Parse(policyReporter)
	if err != nil {
		log.Println(err)
		return
	}

	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiConfig := conf.APIConfig
	authConfig := apiConfig.BasicAuth

	if authConfig.SecretRef != "" && secretClient != nil {
		v, err := secretClient.Get(context.Background(), authConfig.SecretRef)
		if err != nil {
			logger.Error("failed to read secret", zap.String("secret", authConfig.SecretRef), zap.Error(err))
		} else {
			authConfig.Username = v.Username
			authConfig.Password = v.Password
		}
	}

	if development {
		apiRouter.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")

				next.ServeHTTP(w, r)
			})
		})
		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")

				next.ServeHTTP(w, r)
			})
		})
	}

	apiRouter.HandleFunc("/push", api.PushResultHandler(store)).Methods("POST")
	apiRouter.HandleFunc("/result-log", api.ResultHandler(store)).Methods("GET")
	apiRouter.PathPrefix("/v1").Handler(http.StripPrefix("/api", proxy.New(backend, "", false, overwriteHost, apiConfig.Logging, authConfig.Username, authConfig.Password))).Methods("GET")

	for _, c := range conf.APIs {
		if c.SecretRef != "" && secretClient != nil {
			values, err := secretClient.Get(context.Background(), c.SecretRef)
			if err != nil {
				logger.Error("failed to read secret", zap.String("secret", c.SecretRef), zap.Error(err))
			} else {
				c = c.FromValues(values)
			}
		}

		cluster := config.Cluster{
			Name:    c.Name,
			ID:      slug.Make(c.Name),
			Kyverno: len(c.KyvernoAPI) > 0,
		}

		core, err := url.Parse(c.CoreAPI)
		if err != nil {
			logger.Error("failed to configure core proxy", zap.String("name", c.Name), zap.Error(err))
			continue
		}

		apiRouter.
			PathPrefix(fmt.Sprintf("/%s/v1", cluster.ID)).
			Handler(http.StripPrefix(fmt.Sprintf("/api/%s", cluster.ID), proxy.New(core, c.Certificate, c.SkipTLS, overwriteHost, apiConfig.Logging, c.BasicAuth.Username, c.BasicAuth.Password))).Methods("GET")

		logger.Info("core proxy configured", zap.String("name", c.Name))

		if cluster.Kyverno {
			conf.Plugins = AddIfNotExist(conf.Plugins, "kyverno")

			kyverno, err := url.Parse(c.KyvernoAPI)
			if err != nil {
				logger.Error("failed to configure kyverno proxy", zap.String("name", c.Name), zap.Error(err))
				continue
			}

			apiRouter.
				PathPrefix(fmt.Sprintf("/%s/kyverno", cluster.ID)).
				Handler(http.StripPrefix(fmt.Sprintf("/api/%s/kyverno", cluster.ID), proxy.New(kyverno, c.Certificate, c.SkipTLS, overwriteHost, apiConfig.Logging, c.BasicAuth.Username, c.BasicAuth.Password))).Methods("GET")

			logger.Info("kyverno proxy configured", zap.String("name", c.Name))
		}

		conf.Clusters = append(conf.Clusters, cluster)
	}

	if kyvernoPlugin != "" {
		conf.Plugins = AddIfNotExist(conf.Plugins, "kyverno")

		kyverno, err := url.Parse(kyvernoPlugin)
		if err != nil {
			log.Println(err)
			return
		}
		kyvernoProxy := proxy.New(kyverno, "", false, overwriteHost, apiConfig.Logging, apiConfig.BasicAuth.Username, apiConfig.BasicAuth.Password)

		apiRouter.PathPrefix("/kyverno").Handler(http.StripPrefix("/api/kyverno", kyvernoProxy)).Methods("GET")

		logger.Info("kyverno proxy configured")
	}

	apiRouter.HandleFunc("/config", api.ConfigHandler(conf)).Methods("GET")

	if !noUI {
		router.PathPrefix("/").Handler(spaHandler{staticPath: "dist", indexPath: "index.html"}).Methods("GET")
	} else {
		logger.Info("embedded UI disabled")
	}

	if development {
		logger.Info("running", zap.Int("port", port))
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Println(err)
	}
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) && path == h.staticPath {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func AddIfNotExist(list []string, value string) []string {
	for _, v := range list {
		if v == value {
			return list
		}
	}

	return append(list, value)
}

func secretClient(kubeConfig, namespace string) (secrets.Client, error) {
	var k8sConfig *rest.Config
	var err error

	if kubeConfig != "" {
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else {
		k8sConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	clientset, err := k8s.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	return secrets.NewClient(clientset.CoreV1().Secrets(namespace)), nil
}
