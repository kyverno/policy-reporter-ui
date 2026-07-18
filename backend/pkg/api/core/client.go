package core

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kyverno/policy-reporter-ui/pkg/api"
)

type Client struct {
	*api.Client
}

func (c *Client) GetResource(ctx context.Context, id string) (*Resource, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/resource/%s", id), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[Resource](resp.Body)
}

func (c *Client) GetResourceStatusCounts(ctx context.Context, id string, query url.Values) ([]ResourceStatusCount, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/resource/%s/status-counts", id), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[ResourceStatusCount](resp.Body)
}

func (c *Client) GetResourceSeverityCounts(ctx context.Context, id string, query url.Values) ([]ResourceSeverityCount, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/resource/%s/severity-counts", id), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[ResourceSeverityCount](resp.Body)
}

func (c *Client) ListSourceCategoryTree(ctx context.Context, query url.Values) ([]SourceCategoryTree, error) {
	resp, err := c.Get(ctx, "/v2/sources/categories", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[SourceCategoryTree](resp.Body)
}

func (c *Client) ListNamespaceScopedCategories(ctx context.Context, query url.Values) ([]string, error) {
	resp, err := c.Get(ctx, "/v1/namespaced-resources/categories", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[string](resp.Body)
}

func (c *Client) ListClusterScopedCategories(ctx context.Context, query url.Values) ([]string, error) {
	resp, err := c.Get(ctx, "/v1/cluster-resources/categories", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[string](resp.Body)
}

func (c *Client) ListResourceCategories(ctx context.Context, id string, query url.Values) ([]SourceCategoryTree, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/resource/%s/source-categories", id), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[SourceCategoryTree](resp.Body)
}

func (c *Client) GetFindings(ctx context.Context, query url.Values) (*Findings, error) {
	resp, err := c.Get(ctx, "/v2/findings", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[Findings](resp.Body)
}

func (c *Client) GetSeverityFindings(ctx context.Context, query url.Values) (*Findings, error) {
	resp, err := c.Get(ctx, "/v2/severity-findings", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[Findings](resp.Body)
}

func (c *Client) GetNamespaceStatusCounts(ctx context.Context, source string, query url.Values) (NamespaceStatusCounts, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/namespace-scoped/%s/status-counts", source), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeMap[string, map[string]int](resp.Body)
}

func (c *Client) GetNamespaceSeverityCounts(ctx context.Context, source string, query url.Values) (NamespaceStatusCounts, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/namespace-scoped/%s/severity-counts", source), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeMap[string, map[string]int](resp.Body)
}

func (c *Client) GetClusterScopeStatusCounts(ctx context.Context, source string, query url.Values) (map[string]int, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/cluster-scoped/%s/status-counts", source), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeMap[string, int](resp.Body)
}

func (c *Client) GetClusterScopeSeverityCounts(ctx context.Context, source string, query url.Values) (map[string]int, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/cluster-scoped/%s/severity-counts", source), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeMap[string, int](resp.Body)
}

func (c *Client) ListSources(ctx context.Context, query url.Values) ([]string, error) {
	resp, err := c.Get(ctx, "/v2/sources", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[string](resp.Body)
}

func (c *Client) UseResources(ctx context.Context, source string, query url.Values) (bool, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/sources/%s/use-resources", source), query)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	result, err := api.DecodeMap[string, bool](resp.Body)
	if err != nil {
		return false, err
	}

	return result["resources"], nil
}

func (c *Client) ListNamespaces(ctx context.Context, query url.Values) ([]string, error) {
	resp, err := c.Get(ctx, "/v1/namespaces", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[string](resp.Body)
}

func (c *Client) ListNamespacedKinds(ctx context.Context, query url.Values) ([]string, error) {
	resp, err := c.Get(ctx, "/v1/namespaced-resources/kinds", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[string](resp.Body)
}

func (c *Client) ListClusterKinds(ctx context.Context, query url.Values) ([]string, error) {
	resp, err := c.Get(ctx, "/v1/cluster-resources/kinds", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[string](resp.Body)
}

func (c *Client) ListPolicies(ctx context.Context, query url.Values) ([]Policy, error) {
	resp, err := c.Get(ctx, "/v2/policies", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[Policy](resp.Body)
}

func (c *Client) ResolveNamespaceSelector(ctx context.Context, selector map[string]string, query url.Values) ([]string, error) {
	resp, err := c.Post(ctx, "/v2/namespaces/resolve-selector", selector, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[string](resp.Body)
}

func (c *Client) ListNamespaceScopedResourceResults(ctx context.Context, query url.Values) (*Paginated[ResourceResult], error) {
	resp, err := c.Get(ctx, "/v2/namespace-scoped/resource-results", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[Paginated[ResourceResult]](resp.Body)
}

func (c *Client) ListClusterScopedResourceResults(ctx context.Context, query url.Values) (*Paginated[ResourceResult], error) {
	resp, err := c.Get(ctx, "/v2/cluster-scoped/resource-results", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[Paginated[ResourceResult]](resp.Body)
}

func (c *Client) ListResourceResults(ctx context.Context, id string, query url.Values) (*Paginated[PolicyResult], error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/resource/%s/results", id), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[Paginated[PolicyResult]](resp.Body)
}

func (c *Client) ListResourceResourceResults(ctx context.Context, id string, query url.Values) ([]ResourceResult, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v2/resource/%s/resource-results", id), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[ResourceResult](resp.Body)
}

func (c *Client) ListNamespaceScopedResults(ctx context.Context, query url.Values) (*Paginated[PolicyResult], error) {
	resp, err := c.Get(ctx, "/v2/namespace-scoped/results", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[Paginated[PolicyResult]](resp.Body)
}

func (c *Client) ListResultsWithoutResource(ctx context.Context, query url.Values) (*Paginated[PolicyResult], error) {
	resp, err := c.Get(ctx, "/v2/results-without-resources", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[Paginated[PolicyResult]](resp.Body)
}

func (c *Client) ListClusterScopedResults(ctx context.Context, query url.Values) (*Paginated[PolicyResult], error) {
	resp, err := c.Get(ctx, "/v2/cluster-scoped/results", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[Paginated[PolicyResult]](resp.Body)
}

func (c *Client) ListTargets(ctx context.Context) (map[string][]Target, error) {
	resp, err := c.Get(ctx, "/v2/targets", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeMap[string, []Target](resp.Body)
}

func (c *Client) ListTotalResults(ctx context.Context) (*Paginated[ResourceResult], error) {
	resp, err := c.Get(ctx, "/v2/total-results", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[Paginated[ResourceResult]](resp.Body)
}

func New(options []api.ClientOption) (*Client, error) {
	baseClient, err := api.New(options)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Client: baseClient,
	}

	return client, nil
}
