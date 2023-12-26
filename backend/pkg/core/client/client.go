package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
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

func (c *Client) ResolveNamespaceSelector(ctx context.Context, selector map[string]string) ([]string, error) {
	resp, err := c.post(ctx, "/v2/namespaces/resolve-selector", selector)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeList[string](resp.Body)
}

// CreateJSONRequest for the given configuration
func (c *Client) post(ctx context.Context, path string, payload interface{}) (*http.Response, error) {
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

func decodeList[T any](r io.Reader) ([]T, error) {
	list := make([]T, 0)
	err := json.NewDecoder(r).Decode(&list)

	return list, err
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
