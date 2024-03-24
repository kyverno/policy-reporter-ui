package plugin

import (
	"context"
	"fmt"
	"net/url"

	plugin "github.com/kyverno/policy-reporter-plugins/sdk/api"
	"github.com/kyverno/policy-reporter-ui/pkg/api"
)

type Client struct {
	*api.Client
}

func (c *Client) GetPolicies(ctx context.Context) ([]plugin.PolicyListItem, error) {
	resp, err := c.Get(ctx, "/v1/policies", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[plugin.PolicyListItem](resp.Body)
}

func (c *Client) GetPolicy(ctx context.Context, name string) (*plugin.Policy, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v1/policies/%s", name), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[plugin.Policy](resp.Body)
}

func (c *Client) CreateException(ctx context.Context, request *plugin.ExceptionRequest) (*plugin.ExceptionResponse, error) {
	resp, err := c.Post(ctx, "/v1/policies/exception", request)
	if err != nil {
		return nil, fmt.Errorf("create exception endoint failed: %w", err)
	}
	defer resp.Body.Close()

	return api.Decode[plugin.ExceptionResponse](resp.Body)
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
