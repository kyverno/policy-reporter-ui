package client

import (
	"net/http"
	"net/url"
)

type Client interface {
	Get(path string) (*http.Response, error)
}

type httpClient struct {
	client  *http.Client
	baseURL *url.URL
}

func (c *httpClient) Get(path string) (*http.Response, error) {
	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")

	return c.client.Do(req)
}

func NewClient(baseURL string) (Client, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &httpClient{
		baseURL: url,
		client:  &http.Client{},
	}, nil
}
