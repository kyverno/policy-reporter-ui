package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/kyverno/policy-reporter-ui/pkg/api"
	"github.com/kyverno/policy-reporter-ui/pkg/config"
	"github.com/kyverno/policy-reporter-ui/pkg/proxy"
	"github.com/kyverno/policy-reporter-ui/pkg/redis"
	"github.com/kyverno/policy-reporter-ui/pkg/report"
)

var (
	configFile     string
	policyReporter string
	kyvernoPlugin  string
	development    bool
	noUI           bool
	port           int
)

func main() {
	flag.StringVar(&configFile, "config", "./config.yaml", "Path to the config file")
	flag.StringVar(&policyReporter, "policy-reporter", "", "PolicyReporter Host")
	flag.StringVar(&kyvernoPlugin, "kyverno-plugin", "", "Kyverno Plugin Host")
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

	if development {
		log.Println("[INFO] Development Mode enabled")
	}

	router := mux.NewRouter()

	var store report.Store

	if conf.Redis.Enabled {
		if development {
			log.Print("[INFO] Redis Store enabled\n")
		}

		store = redis.NewFromConfig(conf)
	} else {
		if development {
			log.Printf("[INFO] Log Store Size: %d\n", conf.LogSize)
		}

		store = report.NewResultStore(conf.LogSize)
	}

	backend, err := url.Parse(policyReporter)
	if err != nil {
		log.Println(err)
		return
	}

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
	apiRouter.PathPrefix("/v1").Handler(http.StripPrefix("/api", proxy.New(backend, "", false))).Methods("GET")

	for _, c := range conf.APIs {
		cluster := config.Cluster{
			Name:    c.Name,
			ID:      slug.Make(c.Name),
			Kyverno: len(c.KyvernoAPI) > 0,
		}

		core, err := url.Parse(c.CoreAPI)
		if err != nil {
			log.Printf("[ERROR] failed to configure Core Proxy for %s: %s\n", c.Name, err)
			continue
		}

		apiRouter.
			PathPrefix(fmt.Sprintf("/%s/v1", cluster.ID)).
			Handler(http.StripPrefix(fmt.Sprintf("/api/%s", cluster.ID), proxy.New(core, c.Certificate, c.SkipTSL))).Methods("GET")

		log.Printf("[INFO] Core Proxy for %s configured\n", c.Name)

		if cluster.Kyverno {
			conf.Plugins = AddIfNotExist(conf.Plugins, "kyverno")

			kyverno, err := url.Parse(c.KyvernoAPI)
			if err != nil {
				log.Printf("[ERROR] failed to configure Kyverno Proxy for %s: %s\n", c.Name, err)
				continue
			}

			apiRouter.
				PathPrefix(fmt.Sprintf("/%s/kyverno", cluster.ID)).
				Handler(http.StripPrefix(fmt.Sprintf("/api/%s/kyverno", cluster.ID), proxy.New(kyverno, c.Certificate, c.SkipTSL))).Methods("GET")

			log.Printf("[INFO] Kyverno Proxy for %s configured\n", c.Name)
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
		kyvernoProxy := proxy.New(kyverno, "", false)

		apiRouter.PathPrefix("/kyverno").Handler(http.StripPrefix("/api/kyverno", kyvernoProxy)).Methods("GET")

		log.Println("[INFO] Kyverno Plugin Proxy configured")
	}

	apiRouter.HandleFunc("/config", api.ConfigHandler(conf)).Methods("GET")

	if !noUI {
		router.PathPrefix("/").Handler(spaHandler{staticPath: "dist", indexPath: "index.html"}).Methods("GET")
	} else {
		log.Printf("[INFO] Embedded UI disabled")
	}

	if development {
		log.Printf("[INFO] Running on Port %d\n", port)
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
