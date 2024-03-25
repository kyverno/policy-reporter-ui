package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"sync"

	pluginAPI "github.com/kyverno/policy-reporter-plugins/sdk/api"
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/model"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNoClient = errors.New("client for cluster not found")
)

type Service struct {
	endpoints map[string]*model.Endpoints
	configs   map[string]model.SourceConfig
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

	query.Set("sources", source)
	query.Set("policies", policy)

	config := s.configs[source]

	g := &errgroup.Group{}

	var details *pluginAPI.Policy
	if plugin, ok := s.plugin(cluster, source); ok {
		g.Go(func() error {
			details, err = plugin.GetPolicy(ctx, policy)
			zap.L().Error(
				"failed to load policy details from plugin",
				zap.String("cluster", cluster),
				zap.String("source", source),
				zap.String("policy", policy),
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
		Exceptions: s.configs[source].Exceptions,
		Chart: PolicyCharts{
			Findings:       MapFindingsToSourceStatusChart(title, findings),
			NamespaceScope: MapNamespaceScopeChartVariant(title, result, config.EnabledResults()),
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

		status := s.configs[source.Name].EnabledResults()

		list = append(list, Source{
			Name:       source.Name,
			Title:      title,
			Status:     status,
			Categories: categories,
			Chart:      MapCategoriesToChart(title, source.Categories, status),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Title < list[j].Title
	})

	return list, nil
}

func (s *Service) CreateException(ctx context.Context, req ExceptionRequest) (*pluginAPI.ExceptionResponse, error) {
	plugin, ok := s.plugin(req.Cluster, req.Source)
	if !ok {
		return nil, ErrNoClient
	}
	client, err := s.core(req.Cluster)
	if err != nil {
		return nil, err
	}

	resource, err := client.GetResource(ctx, req.Resource)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	var list *core.Paginated[core.PolicyResult]
	if len(req.Policies) == 0 {
		if resource.Namespace != "" {
			list, err = client.ListNamespaceScopedResults(ctx, url.Values{
				"resource_id": []string{req.Resource},
				"status":      []string{StatusFail, StatusWarn},
				"sources":     []string{req.Source},
			})
		} else {
			list, err = client.ListClusterScopedResults(ctx, url.Values{
				"resource_id": []string{req.Resource},
				"status":      []string{StatusFail, StatusWarn},
				"sources":     []string{req.Source},
			})
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get resource results: %w", err)
		}

		results := make(map[string][]string, 0)
		for _, r := range list.Items {
			if _, ok := results[r.Policy]; ok {
				results[r.Policy] = append(results[r.Policy], r.Rule)
			} else {
				results[r.Policy] = []string{r.Rule}
			}
		}

		req.Policies = make([]ExceptionPolicy, 0, len(results))
		for p, rules := range results {
			req.Policies = append(req.Policies, ExceptionPolicy{
				Name:  p,
				Rules: rules,
			})
		}
	}

	request := &pluginAPI.ExceptionRequest{
		Resource: pluginAPI.Resource{
			Name:       resource.Name,
			Namespace:  resource.Namespace,
			Kind:       resource.Kind,
			APIVersion: resource.APIVersion,
		},
		Policies: utils.Map(req.Policies, func(p ExceptionPolicy) *pluginAPI.ExceptionPolicy {
			return &pluginAPI.ExceptionPolicy{
				Name:  p.Name,
				Rules: p.Rules,
			}
		}),
	}

	return plugin.CreateException(ctx, request)
}
func (s *Service) ResourceDetails(ctx context.Context, cluster, id string, query url.Values) (*ResourceDetails, error) {
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

	statusMap := map[string]bool{}
	status := make([]string, 0, 5)

	list := make([]Source, 0, len(sourcesTree))
	for _, source := range sourcesTree {
		categories := make([]string, 0, len(source.Categories))
		for _, category := range source.Categories {
			categories = append(categories, category.Name)
		}

		title := utils.Title(source.Name)

		config := s.configs[source.Name]
		for _, r := range config.EnabledResults() {
			statusMap[r] = true
		}

		list = append(list, Source{
			Name:       source.Name,
			Title:      title,
			Categories: categories,
			Status:     status,
			Exceptions: config.Exceptions,
			Chart:      MapCategoriesToChart(title, source.Categories, config.EnabledResults()),
		})
	}

	for r, ok := range statusMap {
		if ok {
			status = append(status, r)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Title < list[j].Title
	})

	var chart *Chart
	if len(sourcesTree) > 1 {
		chart = MapResourceSourceChart(statusCounts, status)
	}

	return &ResourceDetails{
		Resource: resource,
		Sources:  list,
		Chart:    chart,
		Status:   status,
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

	namespaceResults := make(map[string]core.NamespaceStatusCounts, len(sources))
	clusterResults := make(map[string]map[string]int, len(sources))
	showResults := make([]string, 0, len(sources))

	mx := &sync.Mutex{}
	cmx := &sync.Mutex{}

	statusMap := map[string]bool{}
	status := make([]string, 0, 5)

	var findings *core.Findings
	g.Go(func() error {
		var err error
		findings, err = client.GetFindings(ctx, combinedFilter)

		return err
	})

	for _, source := range sources {
		c := s.configs[source]
		for _, r := range c.EnabledResults() {
			statusMap[r] = true
		}

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

	for r, ok := range statusMap {
		if ok {
			status = append(status, r)
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

	if namespaces == nil {
		namespaces = make([]string, 0)
	}

	singleSource := len(sources) == 1

	var exceptions bool
	if singleSource {
		exceptions = s.configs[sources[0]].Exceptions
	}

	return &Dashboard{
		FilterSources:  make([]string, 0),
		ClusterScope:   clusterScope,
		MultipleSource: len(sources) > 1,
		SingleSource:   singleSource,
		Exceptions:     exceptions,
		Sources:        sources,
		Namespaces:     namespaces,
		ShowResults:    showResults,
		SourcesNavi:    MapFindingSourcesToSourceItem(findings),
		Status:         status,
		Charts: Charts{
			ClusterScope:   clusterResults,
			Findings:       findingChart,
			NamespaceScope: MapNamespaceStatusCountsToCharts(namespaceResults, status),
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

func New(clients map[string]*model.Endpoints, configs map[string]model.SourceConfig) *Service {
	return &Service{clients, configs}
}
