package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fjogeleit/policy-reporter-ui/pkg/client"
)

func PolicyReportHandler(development bool, client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if development {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := client.Get("/policy-reports")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "message": "%s" }`, err.Error())
			return
		}

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "message": "%s" }`, err.Error())
		}
	}
}

func ClusterPolicyReportHandler(development bool, client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if development {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := client.Get("/cluster-policy-reports")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "message": "%s" }`, err.Error())
			return
		}

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "message": "%s" }`, err.Error())
		}
	}
}

func TargetHandler(development bool, client client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if development {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := client.Get("/targets")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "message": "%s" }`, err.Error())
			return
		}

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "message": "%s" }`, err.Error())
		}
	}
}
