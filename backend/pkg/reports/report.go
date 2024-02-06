package reports

import (
	"context"
	"errors"
	"net/url"
	"slices"
	"strings"

	"github.com/kyverno/policy-reporter-plugins/sdk/api"
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/model"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
	"go.uber.org/zap"
)

type ReportGenerator struct {
	endpoints map[string]*model.Endpoints
}

func (g *ReportGenerator) core(cluster string) (*core.Client, bool) {
	endpoints, ok := g.endpoints[cluster]
	if !ok {
		return nil, false
	}

	return endpoints.Core, true
}

func (g *ReportGenerator) plugin(cluster, p string) (*plugin.Client, bool) {
	endpoints, ok := g.endpoints[cluster]
	if !ok {
		return nil, false
	}

	c, ok := endpoints.Plugins[p]

	return c, ok
}

func (g *ReportGenerator) GeneratePerPolicy(ctx context.Context, cluster, source string, filter Filter) ([]*Validation, error) {
	policies, err := g.getPolicies(ctx, cluster, source, filter)
	if err != nil {
		return nil, err
	}

	mapping := make(map[string]*Validation)
	for _, pol := range policies {
		if !filter.IncludesPolicy(pol.Name) {
			continue
		}

		v := &Validation{
			Name:   pol.Title,
			Policy: pol,
			Groups: make(map[string]*Group),
		}

		if pol.Namespace != "" {
			mapping[pol.Namespace+"/"+pol.Name] = v
		}

		mapping[pol.Name] = v
	}

	c, ok := g.core(cluster)
	if !ok {
		return nil, errors.New("core api client not found")
	}

	var namespaces []string

	if len(filter.Namespaces) == 0 {
		namespaces, err = c.ListNamespaces(ctx, url.Values{"sources": []string{source}})
		if err != nil {
			return nil, err
		}
	} else {
		namespaces = filter.Namespaces
	}

	for _, ns := range namespaces {
		results, err := c.ListNamespaceScopedResults(ctx, url.Values{
			"sources":    []string{source},
			"namespaces": []string{ns},
			"categories": filter.Categories,
			"policies":   filter.Policies,
		})
		if err != nil {
			zap.L().Error("failed to list results", zap.Error(err), zap.String("namespace", ns), zap.String("source", source))
			continue
		}

		for _, result := range results.Items {
			val := findResult(mapping, result)
			if val == nil {
				continue
			}

			if result.Status == StatusSkip {
				continue
			}

			_, ok := val.Groups[ns]
			if !ok {
				val.Groups[ns] = &Group{
					Rules:   make(map[string]*Rule),
					Summary: &Summary{},
				}
			}

			rule := result.Rule
			if strings.HasPrefix(rule, "autogen-") {
				rule = strings.TrimPrefix(rule, "autogen-")
			}

			_, ok = val.Groups[ns].Rules[rule]
			if !ok {
				val.Groups[ns].Rules[rule] = &Rule{
					Summary:   &Summary{},
					Resources: make([]*Resource, 0),
				}
			}

			if result.Kind != "" {
				val.Groups[ns].Rules[rule].Resources = append(val.Groups[ns].Rules[rule].Resources, mapResource(result))
			}

			increaseSummary(result.Status, val.Groups[ns].Rules[rule].Summary)
			increaseSummary(result.Status, val.Groups[ns].Summary)
		}
	}

	if filter.ClusterScope {
		results, err := c.ListClusterScopedResults(ctx, url.Values{
			"sources":    []string{source},
			"categories": filter.Categories,
			"policies":   filter.Policies,
		})
		if err != nil {
			zap.L().Error("failed to list results", zap.Error(err), zap.String("source", source))
			return nil, err
		}

		for _, result := range results.Items {
			val := findResult(mapping, result)
			if val == nil {
				continue
			}

			if result.Status == StatusSkip {
				continue
			}

			_, ok := val.Groups[""]
			if !ok {
				val.Groups[""] = &Group{
					Rules:   make(map[string]*Rule),
					Summary: &Summary{},
				}
			}

			rule := result.Rule
			if strings.HasPrefix(rule, "autogen-") {
				rule = strings.TrimPrefix(rule, "autogen-")
			}

			_, ok = val.Groups[""].Rules[rule]
			if !ok {
				val.Groups[""].Rules[rule] = &Rule{
					Summary:   &Summary{},
					Resources: make([]*Resource, 0),
				}
			}

			if result.Kind != "" {
				val.Groups[""].Rules[rule].Resources = append(val.Groups[""].Rules[rule].Resources, mapResource(result))
			}

			increaseSummary(result.Status, val.Groups[""].Rules[rule].Summary)
			increaseSummary(result.Status, val.Groups[""].Summary)
		}
	}

	data := utils.ToList(mapping)

	slices.SortFunc(data, func(a, b *Validation) int {
		if a.Policy.Category != b.Policy.Category {
			return strings.Compare(a.Policy.Category, b.Policy.Category)
		}

		return strings.Compare(a.Name, b.Name)
	})

	return data, nil
}

func (g *ReportGenerator) GeneratePerNamespace(ctx context.Context, cluster, source string, filter Filter) ([]*Validation, error) {
	policies, err := g.getPolicies(ctx, cluster, source, filter)
	if err != nil {
		return nil, err
	}

	groups := make(map[string]api.PolicyListItem)
	for _, pol := range policies {
		if !filter.IncludesPolicy(pol.Name) {
			continue
		}

		if pol.Namespace != "" {
			groups[pol.Name+"/"+pol.Name] = pol
			continue
		}

		groups[pol.Name] = pol
	}

	c, ok := g.core(cluster)
	if !ok {
		return nil, errors.New("core api client not found")
	}

	var namespaces []string

	if len(filter.Namespaces) == 0 {
		namespaces, err = c.ListNamespaces(ctx, url.Values{"sources": []string{source}})
		if err != nil {
			return nil, err
		}
	} else {
		namespaces = filter.Namespaces
	}

	mapping := make(map[string]*Validation)
	for _, ns := range namespaces {
		results, err := c.ListNamespaceScopedResults(ctx, url.Values{
			"sources":    []string{source},
			"namespaces": []string{ns},
			"policies":   filter.Policies,
			"categories": filter.Categories,
		})
		if err != nil {
			zap.L().Error("failed to list results", zap.Error(err), zap.String("namespace", ns), zap.String("source", source))
			continue
		}

		val, ok := mapping[ns]
		if !ok {
			val = &Validation{
				Name:   ns,
				Groups: make(map[string]*Group),
			}

			mapping[ns] = val
		}

		for _, result := range results.Items {
			if result.Status == StatusSkip {
				continue
			}

			_, ok = val.Groups[result.Policy]
			if cache, found := groups[result.Policy]; !ok && found {
				val.Groups[result.Policy] = &Group{
					Name:    result.Policy,
					Policy:  cache,
					Rules:   make(map[string]*Rule),
					Summary: &Summary{},
				}
			} else if !ok {
				val.Groups[result.Policy] = &Group{
					Rules:   make(map[string]*Rule),
					Summary: &Summary{},
				}
			}

			rule := result.Rule
			if strings.HasPrefix(rule, "autogen-") {
				rule = strings.TrimPrefix(rule, "autogen-")
			}

			_, ok = val.Groups[result.Policy].Rules[rule]
			if !ok {
				val.Groups[result.Policy].Rules[rule] = &Rule{
					Summary:   &Summary{},
					Resources: make([]*Resource, 0),
				}
			}

			if result.Kind != "" {
				val.Groups[result.Policy].Rules[rule].Resources = append(val.Groups[result.Policy].Rules[rule].Resources, mapResource(result))
			}

			increaseSummary(result.Status, val.Groups[result.Policy].Rules[rule].Summary)
			increaseSummary(result.Status, val.Groups[result.Policy].Summary)
		}
	}

	data := utils.ToList(mapping)

	slices.SortFunc(data, func(a, b *Validation) int {
		if a.Policy.Category != b.Policy.Category {
			return strings.Compare(a.Policy.Category, b.Policy.Category)
		}

		return strings.Compare(a.Name, b.Name)
	})

	return data, nil
}

func findResult(mapping map[string]*Validation, result core.PolicyResult) *Validation {
	val, ok := mapping[result.Policy]
	if ok {
		return val
	}

	val, ok = mapping[result.Namespace+"/"+result.Policy]
	if ok {
		return val
	}

	return nil
}

func mapResource(result core.PolicyResult) *Resource {
	return &Resource{
		Kind:       result.Kind,
		APIVersion: result.APIVersion,
		Name:       result.Name,
		Status:     result.Status,
	}
}

func increaseSummary(result string, sum *Summary) {
	switch result {
	case StatusPass:
		sum.Pass++
		break
	case StatusWarn:
		sum.Warning++
		break
	case StatusFail:
		sum.Fail++
		break
	case StatusError:
		sum.Error++
		break
	}
}

func (g *ReportGenerator) getPolicies(ctx context.Context, cluster, source string, filter Filter) ([]api.PolicyListItem, error) {
	p, ok := g.plugin(cluster, source)
	if ok {
		policies, err := p.GetPolicies(ctx)
		if err != nil {
			zap.L().Error("failed to load policies from plugin", zap.String("cluster", cluster), zap.String("source", source), zap.Error(err))
		}

		if err == nil {
			if len(filter.Policies) == 0 && len(filter.Categories) == 0 {
				return policies, nil
			}

			return utils.Filter(policies, func(p api.PolicyListItem) bool {
				if len(filter.Policies) > 0 {
					return utils.Contains(filter.Policies, p.Name)
				}

				if len(filter.Categories) > 0 {
					return utils.Contains(filter.Categories, p.Category)
				}

				return true
			}), nil
		}
	}

	c, ok := g.core(cluster)
	if !ok {
		return nil, errors.New("client not found")
	}

	list, err := c.ListPolicies(ctx, url.Values{
		"sources":    []string{source},
		"namespaces": filter.Namespaces,
		"policies":   filter.Policies,
		"categories": filter.Categories,
	})
	if err != nil {
		return nil, err
	}

	return utils.Map(list, func(item core.Policy) api.PolicyListItem {
		return api.PolicyListItem{
			Name:     item.Name,
			Title:    item.Name,
			Category: item.Category,
			Severity: item.Severity,
		}
	}), nil
}
func New(endpoints map[string]*model.Endpoints) *ReportGenerator {
	return &ReportGenerator{endpoints: endpoints}
}
