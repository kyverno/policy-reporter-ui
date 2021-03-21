package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fjogeleit/policy-reporter-ui/pkg/client"
	"github.com/fjogeleit/policy-reporter-ui/pkg/report"
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

func PushResultHandler(store *report.ResultStore) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var result report.Result

		err := json.NewDecoder(req.Body).Decode(&result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "message": "%s" }`, err.Error())
		}

		store.Add(result)
	}
}

func ResultHandler(development bool, store *report.ResultStore) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if development {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		err := json.NewEncoder(w).Encode(store.List())
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
