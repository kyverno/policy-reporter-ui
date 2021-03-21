package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fjogeleit/policy-reporter-ui/pkg/api"
	"github.com/fjogeleit/policy-reporter-ui/pkg/client"
	"github.com/fjogeleit/policy-reporter-ui/pkg/report"
	"github.com/gorilla/mux"
)

var (
	backendURL  string
	development bool
	logSize     int
	port        int
)

func main() {
	flag.StringVar(&backendURL, "backend", "", "PolicyReporter Host")
	flag.IntVar(&port, "port", 8080, "PolicyReporter UI port")
	flag.IntVar(&logSize, "log-size", 200, "Max amount of persistet results")
	flag.BoolVar(&development, "dev", false, "Enable CORS Header for development")
	flag.Parse()

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
