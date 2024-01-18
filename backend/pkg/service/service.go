package service

import (
	"context"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"sync"

	pluginAPI "github.com/kyverno/policy-reporter-plugins/sdk/api"
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNoClient = errors.New("client for cluster not found")
)

type Endpoints struct {
	Core    *core.Client
	Plugins map[string]*plugin.Client
}

type Service struct {
	endpoints map[string]*Endpoints
}

func (s *Service) core(cluster string) (*core.Client, error) {
	endpoints, ok := s.endpoints[cluster]
	if !ok {
		return nil, ErrNoClient
	}

	return endpoints.Core, nil
}

func (s *Service) plugin(cluster, p string) (*plugin.Client, bool) {
	endpoints, ok := s.endpoints[cluster]
	if !ok {
		return nil, false
	}

	c, ok := endpoints.Plugins[p]

	return c, ok
}

func (s *Service) PolicyDetails(ctx context.Context, cluster, source, policy string, query url.Values) (any, error) {
	client, err := s.core(cluster)
	if err != nil {
		return nil, err
	}

	query["sources"] = []string{source}
	query["policies"] = []string{policy}

	g := &errgroup.Group{}

	var details *pluginAPI.Policy
	if plugin, ok := s.plugin(cluster, source); ok {
		g.Go(func() error {
			details, err = plugin.GetPolicy(ctx, policy)
			zap.L().Error(
				"failed to load policy details from plugin",
				zap.String("cluster", cluster),
				zap.String("source", source),
				zap.Error(err),
			)

			return nil
		})
	}

	var namespaces []string
	g.Go(func() error {
		var err error
		namespaces, err = client.ListNamespaces(ctx, query)

		return err
	})

	var findings *core.Findings
	g.Go(func() error {
		var err error
		findings, err = client.GetFindings(ctx, query)

		return err
	})
	var result core.NamespaceStatusCounts
	g.Go(func() error {
		var err error
		result, err = client.GetNamespaceStatusCounts(ctx, source, query)

		return err
	})
	var clusterResult map[string]int
	g.Go(func() error {
		var err error
		clusterResult, err = client.GetClusterStatusCounts(ctx, source, query)

		return err
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}

	title := utils.Title(policy)
	if details != nil {
		title = details.Title
	}

	response := &PolicyDetails{
		Namespaces: namespaces,
		Title:      title,
		Name:       policy,
		Chart: PolicyCharts{
			Findings:       MapFindingsToSourceStatusChart(title, findings),
			NamespaceScope: MapNamespaceStatusCountsToChart(title, result),
			ClusterScope:   clusterResult,
		},
	}

	return MapPolicyDetails(response, details), nil
}

func (s *Service) PolicySources(ctx context.Context, cluster string, query url.Values) ([]Source, error) {
	client, err := s.core(cluster)
	if err != nil {
		return nil, err
	}

	tree, err := client.ListSourceCategoryTree(ctx, query)
	if err != nil {
		return nil, err
	}

	list := make([]Source, 0, len(tree))
	for _, source := range tree {
		categories := make([]string, 0, len(source.Categories))
		for _, category := range source.Categories {
			categories = append(categories, category.Name)
		}

		title := utils.Title(source.Name)

		list = append(list, Source{
			Name:       source.Name,
			Title:      title,
			Categories: categories,
			Chart:      MapCategoriesToChart(title, source.Categories),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Title < list[j].Title
	})

	return list, nil
}

func (s *Service) ResourceDetails(ctx context.Context, cluster string, id string, query url.Values) (*ResourceDetails, error) {
	client, err := s.core(cluster)
	if err != nil {
		return nil, err
	}

	query.Set("resource_id", id)

	g := &errgroup.Group{}

	var resource *core.Resource
	g.Go(func() error {
		var err error
		resource, err = client.GetResource(ctx, id)

		return err
	})

	var statusCounts []core.ResourceStatusCount
	g.Go(func() error {
		var err error
		statusCounts, err = client.GetResourceStatusCounts(ctx, id, query)

		return err
	})

	var sourcesTree []core.SourceCategoryTree
	g.Go(func() error {
		var err error
		sourcesTree, err = client.ListResourceCategories(ctx, id, query)

		return err
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}

	var chart *Chart
	if len(sourcesTree) > 1 {
		chart = MapResourceSourceChart(statusCounts)
	}

	list := make([]Source, 0, len(sourcesTree))
	for _, source := range sourcesTree {
		categories := make([]string, 0, len(source.Categories))
		for _, category := range source.Categories {
			categories = append(categories, category.Name)
		}

		title := utils.Title(source.Name)

		list = append(list, Source{
			Name:       source.Name,
			Title:      title,
			Categories: categories,
			Chart:      MapCategoriesToChart(title, source.Categories),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Title < list[j].Title
	})

	return &ResourceDetails{
		Resource: resource,
		Sources:  list,
		Chart:    chart,
		Results:  SumResourceCounts(statusCounts),
	}, nil
}

func (s *Service) Dashboard(ctx context.Context, cluster string, sources []string, namespaces []string, clusterScope bool, query url.Values) (*Dashboard, error) {
	client, err := s.core(cluster)
	if err != nil {
		return nil, err
	}

	g := &errgroup.Group{}

	combinedFilter, namespaceFilter, clusterFilter := BuildFilters(query)
	combinedFilter.Set("namespaced", strconv.FormatBool(!clusterScope))

	var findings *core.Findings
	g.Go(func() error {
		var err error
		findings, err = client.GetFindings(ctx, combinedFilter)

		return err
	})

	namespaceResults := make(map[string]core.NamespaceStatusCounts, len(sources))
	clusterResults := make(map[string]map[string]int, len(sources))
	showResults := make([]string, 0, len(sources))

	mx := &sync.Mutex{}
	cmx := &sync.Mutex{}

	for _, s := range sources {
		source := s
		g.Go(func() error {
			result, err := client.GetNamespaceStatusCounts(ctx, source, namespaceFilter)
			if err != nil {
				return err
			}

			resources, err := client.UseResources(ctx, source, namespaceFilter)
			if err != nil {
				return err
			}

			mx.Lock()
			namespaceResults[source] = result
			if !resources {
				showResults = append(showResults, source)
			}
			mx.Unlock()

			return nil
		})

		if clusterScope {
			g.Go(func() error {
				result, err := client.GetClusterStatusCounts(ctx, source, clusterFilter)
				if err != nil {
					return err
				}

				cmx.Lock()
				clusterResults[source] = result
				cmx.Unlock()

				return nil
			})
		}
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	var findingChart any
	if len(sources) > 1 {
		findingChart = MapFindingSourcesToFindingCharts(findings)
	} else if len(sources) == 1 {
		findingChart = MapFindingsToSourceStatusChart(sources[0], findings)
	}

	return &Dashboard{
		FilterSources:  make([]string, 0),
		ClusterScope:   clusterScope,
		MultipleSource: len(sources) > 1,
		SingleSource:   len(sources) == 1,
		Sources:        sources,
		Namespaces:     namespaces,
		ShowResults:    showResults,
		SourcesNavi:    MapFindingSourcesToSourceItem(findings),
		Charts: Charts{
			ClusterScope:   clusterResults,
			Findings:       findingChart,
			NamespaceScope: MapNamespaceStatusCountsToCharts(namespaceResults),
		},
		Total: Total{
			Count:     findings.Total,
			PerResult: findings.PerResult,
		},
	}, nil
}

func BuildFilters(baseFilter url.Values) (url.Values, url.Values, url.Values) {
	namespaceFilter := url.Values{}
	clusterFilter := url.Values{}
	combinedFilter := url.Values{}

	for filter, values := range baseFilter {
		if filter == "kinds" || filter == "clusterKinds" {
			continue
		}

		namespaceFilter[filter] = values
		clusterFilter[filter] = values
		combinedFilter[filter] = values
	}

	if val, ok := baseFilter["kinds"]; ok {
		namespaceFilter["kinds"] = val
		combinedFilter["kinds"] = val
	}
	if val, ok := baseFilter["clusterKinds"]; ok {
		clusterFilter["kinds"] = val
		combinedFilter["kinds"] = append(combinedFilter["kinds"], val...)
	}

	return combinedFilter, namespaceFilter, clusterFilter
}

func New(clients map[string]*Endpoints) *Service {
	return &Service{clients}
}
