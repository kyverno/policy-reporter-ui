package api

import (
	"fmt"
	"sort"

	plugin "github.com/kyverno/policy-reporter-plugins/sdk/api"

	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/service"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

func MapSourceCategoryTreeToNavi(sources []core.SourceCategoryTree) []NavigationItem {
	sourceBoards := make([]NavigationItem, 0)
	if len(sources) == 1 {
		for _, category := range sources[0].Categories {
			sourceBoards = append(sourceBoards, NavigationItem{
				Title: utils.Fallback(category.Name, "Other"),
				Path:  fmt.Sprintf("/source/%s/%s", sources[0].Name, category.Name),
			})
		}

		return sourceBoards
	}

	for _, source := range sources {
		count := len(source.Categories)
		if count == 0 {
			sourceBoards = append(sourceBoards, NavigationItem{
				Title: utils.Title(source.Name),
				Path:  fmt.Sprintf("/source/%s", source.Name),
			})
		} else if count == 1 {
			sourceBoards = append(sourceBoards, NavigationItem{
				Title:    utils.Title(source.Name),
				Subtitle: utils.Fallback(source.Categories[0].Name, "Other"),
				Path:     fmt.Sprintf("/source/%s/%s", source.Name, source.Categories[0].Name),
			})
		} else {
			item := NavigationItem{
				Title:    utils.Title(source.Name),
				Path:     fmt.Sprintf("/source/%s", source.Name),
				Children: make([]NavigationItem, 0, len(source.Categories)),
			}

			for _, category := range source.Categories {
				item.Children = append(item.Children, NavigationItem{
					Title: utils.Fallback(category.Name, "Other"),
					Path:  fmt.Sprintf("/source/%s/%s", source.Name, category.Name),
				})
			}

			sourceBoards = append(sourceBoards, item)
		}
	}

	sort.SliceStable(sourceBoards, func(a, b int) bool {
		return sourceBoards[a].Title < sourceBoards[b].Title
	})

	return sourceBoards
}

func MapSourcesToPolicyNavi(sources []core.SourceCategoryTree) []NavigationItem {
	sourceBoards := make([]NavigationItem, 0)
	for _, source := range sources {
		sourceBoards = append(sourceBoards, NavigationItem{
			Title: utils.Title(source.Name),
			Path:  fmt.Sprintf("/policies/%s", source.Name),
		})
	}

	sort.SliceStable(sourceBoards, func(a, b int) bool {
		return sourceBoards[a].Title < sourceBoards[b].Title
	})

	return sourceBoards
}

func MapCustomBoardsToNavi(boards map[string]CustomBoard) []NavigationItem {
	customBoards := make([]NavigationItem, 0, len(boards))
	for _, board := range boards {
		customBoards = append(customBoards, NavigationItem{
			Title: board.Name,
			Path:  fmt.Sprintf("/custom-boards/%s", board.ID),
		})
	}

	sort.SliceStable(customBoards, func(a, b int) bool {
		return customBoards[a].Title < customBoards[b].Title
	})

	return customBoards
}

func MapPoliciesFromCore(policies []core.Policy) map[string][]Policy {
	results := make(map[string][]Policy)
	for _, policy := range policies {
		category := policy.Category
		if category == "" {
			category = "Other"
		}

		if _, ok := results[category]; !ok {
			results[category] = make([]Policy, 0)
		}

		results[category] = append(results[category], Policy{
			Name:     policy.Name,
			Category: policy.Category,
			Severity: policy.Severity,
			Source:   policy.Source,
			Title:    policy.Name,
			Results:  policy.Results,
		})
	}

	return results
}

func MapPluginPolicies(source string, policies []plugin.PolicyListItem, coreList []core.Policy) map[string][]Policy {
	results := make(map[string][]Policy)

	if len(coreList) == 0 {
		return results
	}

	cache := make(map[string]map[string]*core.Policy, len(coreList))
	for _, p := range coreList {
		if _, ok := cache[p.Category]; !ok {
			cache[p.Category] = make(map[string]*core.Policy)
		}
		cache[p.Category][p.Name] = &p
	}

	for _, policy := range policies {
		if _, ok := cache[policy.Category]; !ok {
			cache[policy.Category] = make(map[string]*core.Policy)
		}

		corePolicy := cache[policy.Category][policyID(policy)]
		if corePolicy == nil {
			corePolicy = cache[policy.Category][policy.Name]
		}
		if corePolicy == nil {
			corePolicy = &core.Policy{
				Name:   policy.Name,
				Source: source,
				Results: map[string]int{
					service.StatusPass:  0,
					service.StatusFail:  0,
					service.StatusError: 0,
					service.StatusWarn:  0,
					service.StatusSkip:  0,
				},
			}
		}

		category := policy.Category
		if category == "" {
			category = "Other"
		}

		if _, ok := results[category]; !ok {
			results[category] = make([]Policy, 0)
		}

		results[category] = append(results[category], Policy{
			Namespace:   policy.Namespace,
			Name:        corePolicy.Name,
			Category:    category,
			Severity:    policy.Severity,
			Description: policy.Description,
			Source:      corePolicy.Source,
			Title:       policy.Title,
			Results:     corePolicy.Results,
		})
	}

	return results
}

func policyID(p plugin.PolicyListItem) string {
	if p.Namespace == "" {
		return p.Name
	}

	return fmt.Sprintf("%s/%s", p.Namespace, p.Name)
}
