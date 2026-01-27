package config

import (
	"github.com/gosimple/slug"

	"github.com/kyverno/policy-reporter-ui/pkg/auth"
	"github.com/kyverno/policy-reporter-ui/pkg/customboard"
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
			Permissions: auth.Permissions{
				AccessControl: auth.AccessControl(cl.AccessControl),
			},
		})
	}

	current := ""
	if len(clusters) > 0 {
		current = clusters[0].Slug
	}

	oauth := c.OpenIDConnect.Enabled
	if !oauth {
		oauth = c.OAuth.Enabled
	}

	clusterScope := true
	if c.Boards.ClusterScope.Enabled != nil {
		clusterScope = *c.Boards.ClusterScope.Enabled
	}

	return &api.Config{
		Clusters:    clusters,
		Default:     current,
		OAuth:       oauth,
		Banner:      c.UI.Banner,
		DisplayMode: c.UI.DisplayMode,
		Boards: api.Boards{
			ClusterScope: clusterScope,
			Permissions: auth.Permissions{
				AccessControl: auth.AccessControl(c.Boards.AccessControl),
			},
		},
		Logo: api.Logo{
			Disabled: c.UI.Logo.Disabled,
			Path:     c.UI.Logo.Path,
		},
		Sources: utils.Map(c.Sources, func(s Source) api.Source {
			return api.Source{
				Name:       s.Name,
				ViewType:   s.ViewType,
				Exceptions: s.Exceptions,
				Excludes: api.Excludes{
					NamespaceKinds: s.Excludes.NamespaceKinds,
					ClusterKinds:   s.Excludes.ClusterKinds,
					Results:        s.Excludes.Results,
					Severities:     s.Excludes.Severities,
				},
			}
		}),
	}
}

func MapClusterPermissions(c *Config) map[string]auth.Permissions {
	permissions := make(map[string]auth.Permissions, len(c.Clusters))
	for _, cluster := range c.Clusters {
		permissions[slug.Make(cluster.Name)] = auth.Permissions{
			AccessControl: auth.AccessControl(cluster.AccessControl),
		}
	}

	return permissions
}

func MapFilter(f customboard.Filter) api.Includes {
	if f.NamespaceKinds == nil {
		f.NamespaceKinds = make([]string, 0)
	}
	if f.ClusterKinds == nil {
		f.ClusterKinds = make([]string, 0)
	}
	if f.Results == nil {
		f.Results = make([]string, 0)
	}
	if f.Severities == nil {
		f.Severities = make([]string, 0)
	}

	return api.Includes{
		NamespaceKinds: f.NamespaceKinds,
		ClusterKinds:   f.ClusterKinds,
		Results:        f.Results,
		Severities:     f.Severities,
	}
}
