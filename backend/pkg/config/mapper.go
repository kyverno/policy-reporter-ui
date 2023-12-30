package config

import (
	"github.com/gosimple/slug"

	"github.com/kyverno/policy-reporter-ui/pkg/server/api"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

func MapConfig(c *Config) *api.Config {
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

	return &api.Config{
		Clusters: clusters,
		Default:  current,
		OAuth:    c.OAuth.Enabled,
		DefaultFilter: api.DefaultFilter{
			Resources:        c.UI.DefaultFilter.Resources,
			ClusterResources: c.UI.DefaultFilter.ClusterResources,
		},
		Sources: utils.Map(c.Sources, func(s Source) api.Source {
			return api.Source{
				Name: s.Name,
				Excludes: api.Excludes{
					NamespaceKinds: s.Excludes.NamespaceKinds,
					ClusterKinds:   s.Excludes.ClusterKinds,
				},
			}
		}),
	}
}

func MapCustomBoards(customBoards []CustomBoard) map[string]api.CustomBoard {
	configs := make(map[string]api.CustomBoard, len(customBoards))

	for _, c := range customBoards {
		id := slug.Make(c.Name)

		configs[id] = api.CustomBoard{
			Name: c.Name,
			ID:   id,
			Namespaces: api.Namespaces{
				Selector: c.Namespaces.Selector,
				List:     c.Namespaces.List,
			},
			Sources: api.Sources{
				List: c.Sources.List,
			},
			PolicyReports: api.PolicyReports{
				Selector: c.PolicyReports.Selector,
			}}
	}

	return configs
}
