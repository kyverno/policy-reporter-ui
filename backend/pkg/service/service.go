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
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/model"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

var ErrNoClient = errors.New("client for cluster not found")

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

func (s *Service) Namespaces(ctx context.Context, cluster string, query url.Values) ([]string, error) {
	client, err := s.core(cluster)
	if err != nil {
		return nil, err
	}

	if len(query["sources"]) == 1 {
		config := s.configs[query["sources"][0]]

		query["status"] = config.EnabledResults()
	}

	return client.ListNamespaces(ctx, query)
}

func (s *Service) PolicyDetails(ctx context.Context, cluster, source, policy string, query url.Values) (any, error) {
	client, err := s.core(cluster)
	if err != nil {
		return nil, err
	}

	query.Set("sources", source)
	query.Set("policies", policy)

	config, ok := s.configs[source]
	if ok && len(config.Results) > 0 {
		query["status"] = config.EnabledResults()
	}

	g := &errgroup.Group{}

	var details *pluginAPI.Policy
	if plugin, ok := s.plugin(cluster, source); ok {
		g.Go(func() error {
			details, err = plugin.GetPolicy(ctx, policy)
			if err != nil {
				zap.L().Error(
					"failed to load policy details from plugin",
					zap.String("cluster", cluster),
					zap.String("source", source),
					zap.String("policy", policy),
					zap.Error(err),
				)
			}

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
		Status:     config.EnabledResults(),
		Exceptions: s.configs[source].Exceptions,
		Chart: PolicyCharts{
			Findings:       MapFindingsToSourceStatusChart(title, findings),
			NamespaceScope: MapNamespaceScopeChartVariant(title, result, "", config.EnabledResults(), allStatus),
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

		var chart *Chart
		if s.configs[source.Name].ViewType == model.Severity {
			chart = MapCategorySeveritiesToChart(title, source.Categories, []string{})
			status = []string{"summary"}
		} else {
			chart = MapCategoryStatusToChart(title, source.Categories, status)
		}

		list = append(list, Source{
			Name:       source.Name,
			Title:      title,
			Status:     status,
			Categories: categories,
			Chart:      chart,
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
		var categories []string
		if req.Category != "" {
			categories = []string{req.Category}
		}

		if resource.Namespace != "" {
			list, err = client.ListNamespaceScopedResults(ctx, url.Values{
				"resource_id": []string{req.Resource},
				"status":      []string{StatusFail, StatusWarn},
				"sources":     []string{req.Source},
				"categories":  categories,
			})
		} else {
			list, err = client.ListClusterScopedResults(ctx, url.Values{
				"resource_id": []string{req.Resource},
				"status":      []string{StatusFail, StatusWarn},
				"sources":     []string{req.Source},
				"categories":  categories,
			})
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get resource results: %w", err)
		}

		results := make(map[string][]ExceptionRule, 0)
		for _, r := range list.Items {
			if _, ok := results[r.Policy]; ok {
				results[r.Policy] = append(results[r.Policy], ExceptionRule{Name: r.Rule, Props: r.Properties})
			} else {
				results[r.Policy] = []ExceptionRule{{Name: r.Rule, Props: r.Properties}}
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
				Name: p.Name,
				Rules: utils.Map(p.Rules, func(rule ExceptionRule) pluginAPI.ExceptionRule {
					return pluginAPI.ExceptionRule{
						Name:  rule.Name,
						Props: rule.Props,
					}
				}),
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

	var severityCounts []core.ResourceSeverityCount
	g.Go(func() error {
		var err error
		severityCounts, err = client.GetResourceSeverityCounts(ctx, id, query)

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

	severityMap := map[string]bool{}
	severities := make([]string, 0, 6)

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
		for _, r := range config.EnabledSeverities() {
			severityMap[r] = true
		}

		var plugin bool
		if _, ok := s.plugin(cluster, source.Name); ok {
			plugin = true
		}

		list = append(list, Source{
			Name:       source.Name,
			Title:      title,
			Categories: categories,
			Status:     status,
			Exceptions: config.Exceptions,
			Plugin:     plugin,
			Chart:      MapCategoryStatusToChart(title, source.Categories, config.EnabledResults()),
		})
	}

	for r, ok := range statusMap {
		if ok {
			status = append(status, r)
		}
	}
	for r, ok := range severityMap {
		if ok {
			severities = append(severities, r)
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
		Resource:        resource,
		Sources:         list,
		Chart:           chart,
		Status:          status,
		Results:         SumResourceCounts(statusCounts),
		SeverityResults: SumResourceSeverityCounts(severityCounts),
		Severities:      severities,
	}, nil
}

func (s *Service) Dashboard(ctx context.Context, o DashboardOptions, query url.Values) (*Dashboard, error) {
	if s.viewType(o.Sources) == model.Severity {
		config, ok := s.configs[o.Sources[0]]
		if ok && config.ViewType == model.Severity {
			return s.SeverityDashboard(ctx, o, query)
		}
	}

	client, err := s.core(o.Cluster)
	if err != nil {
		return nil, err
	}

	g := &errgroup.Group{}

	combinedFilter, namespaceFilter, clusterFilter := BuildFilters(query)
	combinedFilter.Set("namespaced", strconv.FormatBool(!o.ClusterScope))

	namespaceResults := make(map[string]core.NamespaceStatusCounts, len(o.Sources))
	clusterResults := make(map[string]map[string]int, len(o.Sources))
	showResults := make([]string, 0, len(o.Sources))

	mx := &sync.Mutex{}
	cmx := &sync.Mutex{}

	if len(o.Status) == 0 {
		o.Status = s.filterEnabled(o.Sources, func(c model.SourceConfig) []string {
			return c.EnabledResults()
		})
	}

	if len(o.Severities) == 0 {
		o.Severities = s.filterEnabled(o.Sources, func(c model.SourceConfig) []string {
			return c.EnabledSeverities()
		})
	}

	combinedFilter["status"] = o.Status
	namespaceFilter["status"] = o.Status
	clusterFilter["status"] = o.Status

	var findings *core.Findings
	g.Go(func() error {
		var err error
		findings, err = client.GetFindings(ctx, combinedFilter)

		return err
	})

	for _, source := range o.Sources {
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

		if o.ClusterScope {
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

	if o.Namespaces == nil {
		o.Namespaces = make([]string, 0)
	}

	singleSource := len(o.Sources) == 1

	var exceptions bool
	if singleSource {
		exceptions = s.configs[o.Sources[0]].Exceptions
	}

	var findingChart any
	if len(o.Sources) > 1 {
		findingChart = MapFindingSourcesToFindingCharts(findings)
	} else if len(o.Sources) == 1 {
		findingChart = MapFindingsToSourceStatusChart(o.Sources[0], findings)
	} else {
		findingChart = MapFindingsToSourceStatusChart("", &core.Findings{})
	}

	return &Dashboard{
		Type:           model.Status,
		FilterSources:  make([]string, 0),
		ClusterScope:   o.ClusterScope,
		MultipleSource: len(o.Sources) > 1,
		SingleSource:   singleSource,
		Exceptions:     exceptions,
		Sources:        o.Sources,
		Namespaces:     o.Namespaces,
		ShowResults:    showResults,
		SourcesNavi:    MapFindingSourcesToSourceItem(findings),
		Status:         o.Status,
		Severities:     o.Severities,
		Display:        o.Display,
		NamespaceKinds: o.NamespaceKinds,
		ClusterKinds:   o.ClusterKinds,
		Charts: Charts{
			ClusterScope:   clusterResults,
			Findings:       findingChart,
			NamespaceScope: MapNamespaceStatusCountsToCharts(namespaceResults, model.Status, o.Status, allStatus),
		},
		Total: Total{
			Count:     findings.Total,
			PerResult: findings.PerResult,
		},
	}, nil
}

func (s *Service) SeverityDashboard(ctx context.Context, o DashboardOptions, query url.Values) (*Dashboard, error) {
	client, err := s.core(o.Cluster)
	if err != nil {
		return nil, err
	}

	g := &errgroup.Group{}

	combinedFilter, namespaceFilter, clusterFilter := BuildFilters(query)
	combinedFilter.Set("namespaced", strconv.FormatBool(!o.ClusterScope))

	namespaceResults := make(map[string]core.NamespaceStatusCounts, len(o.Sources))
	clusterResults := make(map[string]map[string]int, len(o.Sources))
	showResults := make([]string, 0, len(o.Sources))

	mx := &sync.Mutex{}
	cmx := &sync.Mutex{}

	if len(o.Status) == 0 {
		o.Status = s.filterEnabled(o.Sources, func(c model.SourceConfig) []string {
			return c.EnabledResults()
		})
	}

	if len(o.Severities) == 0 {
		o.Severities = s.filterEnabled(o.Sources, func(c model.SourceConfig) []string {
			return c.EnabledSeverities()
		})
	}

	combinedFilter["severity"] = o.Severities
	namespaceFilter["severiy"] = o.Severities
	clusterFilter["severity"] = o.Severities

	var findings *core.Findings
	g.Go(func() error {
		var err error
		findings, err = client.GetSeverityFindings(ctx, combinedFilter)

		return err
	})

	for _, source := range o.Sources {
		g.Go(func() error {
			result, err := client.GetNamespaceSeverityCounts(ctx, source, namespaceFilter)
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

		if o.ClusterScope {
			g.Go(func() error {
				result, err := client.GetClusterSeverityCounts(ctx, source, clusterFilter)
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

	if o.Namespaces == nil {
		o.Namespaces = make([]string, 0)
	}

	singleSource := len(o.Sources) == 1

	var exceptions bool
	if singleSource {
		exceptions = s.configs[o.Sources[0]].Exceptions
	}

	var findingChart any
	if len(o.Sources) > 1 {
		findingChart = MapFindingSourcesToFindingCharts(findings)
	} else if len(o.Sources) == 1 {
		findingChart = MapSeverityFindingsToSourceStatusChart(o.Sources[0], findings)
	} else {
		findingChart = MapSeverityFindingsToSourceStatusChart("", &core.Findings{})
	}

	return &Dashboard{
		Type:           model.Severity,
		FilterSources:  make([]string, 0),
		ClusterScope:   o.ClusterScope,
		MultipleSource: len(o.Sources) > 1,
		SingleSource:   singleSource,
		Exceptions:     exceptions,
		Sources:        o.Sources,
		Namespaces:     o.Namespaces,
		ShowResults:    showResults,
		SourcesNavi:    MapFindingSourcesToSourceItem(findings),
		Status:         o.Status,
		Severities:     o.Severities,
		Display:        o.Display,
		NamespaceKinds: o.NamespaceKinds,
		ClusterKinds:   o.ClusterKinds,
		Charts: Charts{
			ClusterScope:   clusterResults,
			Findings:       findingChart,
			NamespaceScope: MapNamespaceStatusCountsToCharts(namespaceResults, model.Severity, o.Severities, allSeverities),
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

func (s *Service) filterEnabled(sources []string, call func(c model.SourceConfig) []string) []string {
	mapping := map[string]bool{}
	list := make([]string, 0)

	for _, source := range sources {
		config := s.configs[source]
		for _, r := range call(config) {
			mapping[r] = true
		}
	}
	for r, ok := range mapping {
		if ok {
			list = append(list, r)
		}
	}

	return list
}

func (s *Service) viewType(sources []string) string {
	if len(sources) == 0 {
		return model.Status
	}

	for _, source := range sources {
		config := s.configs[source]
		if config.ViewType != model.Severity {
			return model.Status
		}
	}

	return model.Severity
}

func New(clients map[string]*model.Endpoints, configs map[string]model.SourceConfig) *Service {
	return &Service{clients, configs}
}
