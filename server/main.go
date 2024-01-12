package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"go.uber.org/zap"

	"github.com/kyverno/policy-reporter-ui/pkg/api"
	"github.com/kyverno/policy-reporter-ui/pkg/config"
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

	if overwriteHost {
		conf.APIConfig.OverwriteHost = overwriteHost
	}

	resolver := config.NewResolver(conf, kubeConfig, development)
	resolver.InitClusters(kyvernoPlugin != "")
	logger := resolver.Logger()

	if err := resolver.InitSecretClient(); err != nil {
		logger.Warn("unable to setup secret reader, secretRefs can't be resolved", zap.Error(err))
	}

	if development {
		logger.Info("Development Mode enabled")
	}

	router := mux.NewRouter()
	store := resolver.Store()

	apiRouter := router.PathPrefix("/api/").Subrouter()

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

	coreProxy, kyvernoProxy, err := resolver.InternalProxies(context.Background(), policyReporter, kyvernoPlugin)
	if err != nil {
		logger.Error("unable to setup proxies", zap.Error(err), zap.String("api", policyReporter), zap.String("plugin", kyvernoPlugin))
		return
	}

	apiRouter.PathPrefix("/v1").Handler(http.StripPrefix("/api", coreProxy)).Methods("GET")
	
	if kyvernoProxy != nil {
		conf.Plugins = AddIfNotExist(conf.Plugins, "kyverno")

		apiRouter.PathPrefix("/kyverno").Handler(http.StripPrefix("/api/kyverno", kyvernoProxy)).Methods("GET")

		logger.Info("kyverno proxy configured")
	}

	for _, c := range conf.APIs {
		cluster := config.Cluster{
			Name:    c.Name,
			ID:      slug.Make(c.Name),
			Kyverno: c.KyvernoAPI != "",
		}

		core, plugin, err := resolver.ExternalProxies(context.Background(), c)
		if err != nil {
			logger.Error("failed to configure api proxies", zap.String("name", c.Name), zap.Error(err))
			continue
		}

		apiRouter.
			PathPrefix(fmt.Sprintf("/%s/v1", cluster.ID)).
			Handler(http.StripPrefix(fmt.Sprintf("/api/%s", cluster.ID), core)).Methods("GET")

		logger.Info("core proxy configured", zap.String("name", c.Name))

		if plugin != nil {
			conf.Plugins = AddIfNotExist(conf.Plugins, "kyverno")

			apiRouter.
				PathPrefix(fmt.Sprintf("/%s/kyverno", cluster.ID)).
				Handler(http.StripPrefix(fmt.Sprintf("/api/%s/kyverno", cluster.ID), plugin)).Methods("GET")

			logger.Info("kyverno proxy configured", zap.String("name", c.Name))
		}

		conf.Clusters = append(conf.Clusters, cluster)
	}

	apiRouter.HandleFunc("/config", api.ConfigHandler(conf)).Methods("GET")

	if !noUI {
		router.PathPrefix("/").Handler(spaHandler{staticPath: "dist", indexPath: "index.html"}).Methods("GET")
	} else {
		logger.Info("embedded UI disabled")
	}

	logger.Info("server starting", zap.Int("port", port))

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
