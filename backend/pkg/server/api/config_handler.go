package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefaultFilter struct {
	Resources        []string `json:"resources"`
	ClusterResources []string `json:"clusterResources"`
}

type Excludes struct {
	NamespaceKinds []string `json:"namespaceKinds"`
	ClusterKinds   []string `json:"clusterKinds"`
}

type Source struct {
	Name     string   `json:"name"`
	Excludes Excludes `json:"excludes"`
}

type Cluster struct {
	Name    string   `json:"name"`
	Slug    string   `json:"slug"`
	Plugins []string `json:"plugins"`
}

type Config struct {
	Clusters      []Cluster     `json:"clusters"`
	Sources       []Source      `json:"sources"`
	Default       string        `json:"default"`
	DefaultFilter DefaultFilter `json:"defaultFilter"`
}

func ConfigHandler(conf Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, conf)
	}
}
