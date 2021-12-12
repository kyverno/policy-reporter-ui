package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/kyverno/policy-reporter-ui/pkg/api"
	"github.com/kyverno/policy-reporter-ui/pkg/config"
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

	if development {
		log.Println("[INFO] Development Mode enabled")
	}
	if development {
		log.Printf("[INFO] Log Store Size: %d\n", conf.LogSize)
	}

	router := mux.NewRouter()

	store := report.NewResultStore(conf.LogSize)

	backend, err := url.Parse(policyReporter)
	if err != nil {
		log.Println(err)
		return
	}

	coreProxy := httputil.NewSingleHostReverseProxy(backend)

	apiRouter := router.PathPrefix("/api/").Subrouter()

	if development {
		apiRouter.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")

				next.ServeHTTP(w, r)
			})
		})
	}

	apiRouter.HandleFunc("/push", api.PushResultHandler(store)).Methods("POST")
	apiRouter.HandleFunc("/result-log", api.ResultHandler(store)).Methods("GET")
	apiRouter.PathPrefix("/v1").Handler(http.StripPrefix("/api", coreProxy)).Methods("GET")

	if kyvernoPlugin != "" {
		conf.Plugins = append(conf.Plugins, "kyverno")

		kyverno, err := url.Parse(kyvernoPlugin)
		if err != nil {
			log.Println(err)
			return
		}
		kyvernoProxy := httputil.NewSingleHostReverseProxy(kyverno)

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
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
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
