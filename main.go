package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kyverno/policy-reporter-ui/pkg/api"
	"github.com/kyverno/policy-reporter-ui/pkg/client"
	"github.com/kyverno/policy-reporter-ui/pkg/report"
)

var (
	backendURL    string
	kyvernoPlugin string
	development   bool
	logSize       int
	port          int
)

func main() {
	flag.StringVar(&backendURL, "backend", "", "PolicyReporter Host")
	flag.StringVar(&kyvernoPlugin, "kyverno-plugin", "", "Kyverno Plugin Host")
	flag.IntVar(&port, "port", 8080, "PolicyReporter UI port")
	flag.IntVar(&logSize, "log-size", 200, "Max amount of persistet results")
	flag.BoolVar(&development, "dev", false, "Enable CORS Header for development")
	flag.Parse()

	var plugins []string

	if development {
		log.Println("[INFO] Development Mode enabled")
	}
	if development {
		log.Printf("[INFO] Log Store Size: %d\n", logSize)
	}

	c, err := client.NewClient(backendURL)
	if err != nil {
		log.Println(err)
	}

	router := mux.NewRouter()

	store := report.NewResultStore(logSize)

	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.HandleFunc("/push", api.PushResultHandler(store)).Methods("POST")
	apiRouter.HandleFunc("/result-log", api.ResultHandler(development, store)).Methods("GET")
	apiRouter.HandleFunc("/targets", api.TargetHandler(development, c)).Methods("GET")
	apiRouter.HandleFunc("/policy-reports", api.PolicyReportHandler(development, c)).Methods("GET")
	apiRouter.HandleFunc("/cluster-policy-reports", api.ClusterPolicyReportHandler(development, c)).Methods("GET")

	if kyvernoPlugin != "" {
		plugins = append(plugins, "kyverno")

		kC, err := client.NewClient(kyvernoPlugin)
		if err != nil {
			log.Println(err)
		}

		apiRouter.HandleFunc("/kyverno/policies", api.KyvernoPolicyHandler(development, kC)).Methods("GET")
	}

	apiRouter.HandleFunc("/plugins", api.PluginHandler(development, plugins)).Methods("GET")

	handler := http.FileServer(http.Dir("dist"))

	router.PathPrefix("/js").Handler(handler).Methods("GET")
	router.PathPrefix("/css").Handler(handler).Methods("GET")
	router.PathPrefix("/favicon.ico").Handler(handler).Methods("GET")
	router.PathPrefix("/").Handler(IndexHandler("/dist/index.html")).Methods("GET")

	if development {
		log.Printf("[INFO] Running on Port %d\n", port)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Println(err)
	}
}

func IndexHandler(path string) http.Handler {
	f, err := http.Dir(".").Open(path)
	if err != nil {
		log.Println(err)
	}
	d, err := f.Stat()
	if err != nil {
		log.Println(err)
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, d.Name(), d.ModTime(), f.(io.ReadSeeker))
	}

	return http.HandlerFunc(fn)
}
