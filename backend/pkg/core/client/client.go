package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type BasicAuth struct {
	Username string
	Password string
}

type Client struct {
	baseURL string
	http    *http.Client
	auth    *BasicAuth
}

func (c *Client) GetResource(ctx context.Context, id string) (*Resource, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/v2/resource/%s", id), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decode[Resource](resp.Body)
}
func (c *Client) GetResourceStatusCounts(ctx context.Context, id string, query url.Values) ([]ResourceStatusCount, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/v2/resource/%s/status-counts", id), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeList[ResourceStatusCount](resp.Body)
}

func (c *Client) ListSourceCategoryTree(ctx context.Context, query url.Values) ([]SourceCategoryTree, error) {
	resp, err := c.get(ctx, "/v2/sources/categories", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeList[SourceCategoryTree](resp.Body)
}

func (c *Client) GetFindings(ctx context.Context, query url.Values) (*Findings, error) {
	resp, err := c.get(ctx, "/v2/findings", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decode[Findings](resp.Body)
}

func (c *Client) GetNamespaceStatusCounts(ctx context.Context, source string, query url.Values) (NamespaceStatusCounts, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/v2/namespace-scoped/%s/status-counts", source), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeMap[string, map[string]int](resp.Body)
}

func (c *Client) GetClusterStatusCounts(ctx context.Context, source string, query url.Values) (map[string]int, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/v2/cluster-scoped/%s/status-counts", source), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeMap[string, int](resp.Body)
}

func (c *Client) ListSources(ctx context.Context, query url.Values) ([]string, error) {
	resp, err := c.get(ctx, "/v2/sources", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeList[string](resp.Body)
}

func (c *Client) ListNamespaces(ctx context.Context, query url.Values) ([]string, error) {
	resp, err := c.get(ctx, "/v1/namespaces", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeList[string](resp.Body)
}

func (c *Client) ResolveNamespaceSelector(ctx context.Context, selector map[string]string) ([]string, error) {
	resp, err := c.post(ctx, "/v2/namespaces/resolve-selector", selector)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeList[string](resp.Body)
}

// CreateJSONRequest for the given configuration
func (c *Client) post(ctx context.Context, path string, payload any) (*http.Response, error) {
	body := new(bytes.Buffer)

	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	if c.auth != nil {
		req.SetBasicAuth(c.auth.Username, c.auth.Password)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", "Policy Reporter UI")

	return c.http.Do(req)
}

// CreateJSONRequest for the given configuration
func (c *Client) get(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	if c.auth != nil {
		req.SetBasicAuth(c.auth.Username, c.auth.Password)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", "Policy Reporter UI")
	req.URL.RawQuery = query.Encode()

	return c.http.Do(req)
}

func decodeList[T any](r io.Reader) ([]T, error) {
	list := make([]T, 0)
	err := json.NewDecoder(r).Decode(&list)

	return list, err
}

func decodeMap[R comparable, T any](r io.Reader) (map[R]T, error) {
	mapping := make(map[R]T)
	err := json.NewDecoder(r).Decode(&mapping)

	return mapping, err
}

func decode[T any](r io.Reader) (*T, error) {
	model := new(T)
	err := json.NewDecoder(r).Decode(model)

	return model, err
}

func New(options []ClientOption) (*Client, error) {
	client := &Client{
		http: newHTTPClient(),
	}

	for _, o := range options {
		if err := o(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{},
		},
		Timeout: 10 * time.Second,
	}
}
