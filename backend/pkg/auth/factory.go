package auth

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/yandex"

	"github.com/kyverno/policy-reporter-ui/pkg/api"
)

type Options struct {
	Scopes       []string
	DiscoveryURL string
}

func NewProvider(provider string, clientKey string, secret string, callbackURL string, scopes []string) goth.Provider {
	switch provider {
	case "amazon":
		p := amazon.New(clientKey, secret, callbackURL, scopes...)
		p.HTTPClient = api.NewHTTPClient()

		return p
	case "gitlab":
		p := gitlab.New(clientKey, secret, callbackURL, scopes...)
		p.HTTPClient = api.NewHTTPClient()

		return p
	case "github":
		p := github.New(clientKey, secret, callbackURL, scopes...)
		p.HTTPClient = api.NewHTTPClient()

		return p
	case "apple":
		return apple.New(clientKey, secret, callbackURL, api.NewHTTPClient(), apple.ScopeName, apple.ScopeEmail)
	case "google":
		p := google.New(clientKey, secret, callbackURL, scopes...)
		p.HTTPClient = api.NewHTTPClient()

		return p
	case "yandex":
		p := yandex.New(clientKey, secret, callbackURL, scopes...)
		p.HTTPClient = api.NewHTTPClient()

		return p
	case "azuread":
		p := azuread.New(clientKey, secret, callbackURL, nil, scopes...)
		p.HTTPClient = api.NewHTTPClient()

		return p
	}

	return nil
}
