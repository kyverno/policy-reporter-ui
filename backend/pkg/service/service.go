package service

import (
	"context"
	"errors"
	"net/url"
	"sort"
	"sync"

	core "github.com/kyverno/policy-reporter-ui/pkg/core/client"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNoClient = errors.New("client for cluster not found")
)

type Service struct {
	clients map[string]*core.Client
}

func (s *Service) PolicyDetails(ctx context.Context, cluster, source, policy string, query url.Values) (any, error) {
	client, ok := s.clients[cluster]
	if !ok {
		return nil, errors.New("cluster not found")
	}

	query["sources"] = []string{source}
	query["policies"] = []string{policy}

	g := &errgroup.Group{}

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

	return &PolicyDetails{
		Namespaces: namespaces,
		Title:      utils.Title(policy),
		Name:       policy,
		Chart: PolicyCharts{
			Findings:       MapFindingsToSourceStatusChart(utils.Title(policy), findings),
			NamespaceScope: MapNamespaceStatusCountsToChart(utils.Title(policy), result),
			ClusterScope:   clusterResult,
		},
	}, nil
}

func (s *Service) PolicySources(ctx context.Context, cluster string, query url.Values) ([]Source, error) {
	client, ok := s.clients[cluster]
	if !ok {
		return nil, errors.New("cluster not found")
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
	client, ok := s.clients[cluster]
	if !ok {
		return nil, errors.New("cluster not found")
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
		sourcesTree, err = client.ListSourceCategoryTree(ctx, query)

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
	client, ok := s.clients[cluster]
	if !ok {
		return nil, errors.New("cluster not found")
	}

	g := &errgroup.Group{}

	combinedFilter, namespaceFilter, clusterFilter := BuildFilters(query)

	var findings *core.Findings
	g.Go(func() error {
		var err error
		findings, err = client.GetFindings(ctx, combinedFilter)

		return err
	})

	namespaceResults := make(map[string]core.NamespaceStatusCounts, len(sources))
	clusterResults := make(map[string]map[string]int, len(sources))

	mx := &sync.Mutex{}
	cmx := &sync.Mutex{}

	for _, s := range sources {
		source := s
		g.Go(func() error {
			result, err := client.GetNamespaceStatusCounts(ctx, source, namespaceFilter)
			if err != nil {
				return err
			}

			mx.Lock()
			namespaceResults[source] = result
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
	} else {
		findingChart = MapFindingsToSourceStatusChart(sources[0], findings)
	}

	return &Dashboard{
		FilterSources:  make([]string, 0),
		ClusterScope:   clusterScope,
		MultipleSource: len(sources) > 1,
		SingleSource:   len(sources) == 1,
		Sources:        sources,
		Namespaces:     namespaces,
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

func New(clients map[string]*core.Client) *Service {
	return &Service{clients}
}
