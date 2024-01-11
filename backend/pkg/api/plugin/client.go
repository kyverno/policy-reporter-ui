package plugin

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kyverno/policy-reporter-ui/pkg/api"
)

type Client struct {
	*api.Client
}

func (c *Client) ListPolicies(ctx context.Context, query url.Values) ([]Policy, error) {
	resp, err := c.Get(ctx, "/policies", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.DecodeList[Policy](resp.Body)
}

func (c *Client) GetPolicy(ctx context.Context, name string) (*PolicyDetails, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/policies/%s", name), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return api.Decode[PolicyDetails](resp.Body)
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
