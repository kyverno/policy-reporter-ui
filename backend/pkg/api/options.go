package api

import (
	"fmt"
	"net/http"
)

type BasicAuth struct {
	Username string
	Password string
}

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
		certs, err := LoadCerts(path)
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
		client.http.Transport = NewLoggingRoundTripper(client.http.Transport)

		return nil
	}
}
