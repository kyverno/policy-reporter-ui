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

	return &api.Config{
		Clusters:    clusters,
		Default:     current,
		OAuth:       oauth,
		Banner:      c.UI.Banner,
		DisplayMode: c.UI.DisplayMode,
		Boards: api.Boards{
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

func MapCustomBoards(customBoards []customboard.CustomBoard) map[string]api.CustomBoard {
	configs := make(map[string]api.CustomBoard, len(customBoards))

	for _, c := range customBoards {
		id := slug.Make(c.Name)

		configs[id] = api.CustomBoard{
			Name:    c.Name,
			ID:      id,
			Display: c.Display,
			Filter:  MapFilter(c.Filter.Include),
			Namespaces: api.Namespaces{
				Selector: c.Namespaces.Selector,
				List:     c.Namespaces.List,
			},
			Sources: api.Sources{
				List: c.Sources.List,
			},
			Permissions: auth.Permissions{
				AccessControl: auth.AccessControl(c.AccessControl),
			},
			PolicyReports: api.PolicyReports{
				Selector: c.PolicyReports.Selector,
			},
			ClusterScope: c.ClusterScope.Enabled,
		}
	}

	return configs
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
