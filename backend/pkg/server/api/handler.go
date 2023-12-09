package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefaultFilter struct {
	Resources        []string `json:"resources"`
	ClusterResources []string `json:"clusterResources"`
}

type Cluster struct {
	Name    string   `json:"name"`
	Slug    string   `json:"slug"`
	Plugins []string `json:"plugins"`
}

type Config struct {
	Clusters      []Cluster     `json:"clusters"`
	Default       string        `json:"default"`
	DefaultFilter DefaultFilter `json:"defaultFilter"`
}

func ConfigHandler(conf Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, conf)
	}
}
