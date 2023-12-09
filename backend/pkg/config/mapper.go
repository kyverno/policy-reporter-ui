package config

import (
	"github.com/gosimple/slug"

	"github.com/kyverno/policy-reporter-ui/pkg/server/api"
)

func MapConfig(c *Config) api.Config {
	clusters := make([]api.Cluster, 0, len(c.Clusters))
	for _, cl := range c.Clusters {
		plugins := make([]string, 0, len(cl.Plugins))
		for _, pl := range cl.Plugins {
			plugins = append(plugins, pl.Name)
		}

		clusters = append(clusters, api.Cluster{
			Name:    cl.Name,
			Slug:    slug.Make(cl.Name),
			Plugins: plugins,
		})
	}

	current := ""
	if len(clusters) > 0 {
		current = clusters[0].Slug
	}

	return api.Config{
		Clusters: clusters,
		Default:  current,
		DefaultFilter: api.DefaultFilter{
			Resources:        c.UI.DefaultFilter.Resources,
			ClusterResources: c.UI.DefaultFilter.ClusterResources,
		},
	}
}
