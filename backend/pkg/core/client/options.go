package client

import (
	"fmt"
	"net/http"

	"github.com/kyverno/policy-reporter-ui/pkg/core/utils"
)

type ClientOption = func(*Client) error

func WithBaseURL(url string) ClientOption {
	return func(client *Client) error {
		client.baseURL = url

		return nil
	}
}

func WithBaseAuth(auth BasicAuth) ClientOption {
	return func(client *Client) error {
		client.auth = &auth

		return nil
	}
}

func WithCertificate(path string) ClientOption {
	return func(client *Client) error {
		certs, err := utils.LoadCerts(path)
		if err != nil {
			return fmt.Errorf("with certificate failed: %w", err)
		}

		client.http.Transport.(*http.Transport).TLSClientConfig.RootCAs = certs

		return nil
	}
}

func WithSkipTLS() ClientOption {
	return func(client *Client) error {
		client.http.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = true

		return nil
	}
}

func WithLogging() ClientOption {
	return func(client *Client) error {
		client.http.Transport = newLoggingRoundTripper(client.http.Transport)

		return nil
	}
}
