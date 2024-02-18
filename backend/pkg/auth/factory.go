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
)

type Options struct {
	Scopes       []string
	DiscoveryURL string
}

func NewProvider(provider string, clientKey string, secret string, callbackURL string, scopes []string) goth.Provider {
	switch provider {
	case "amazon":
		return amazon.New(clientKey, secret, callbackURL, scopes...)
	case "gitlab":
		return gitlab.New(clientKey, secret, callbackURL, scopes...)
	case "github":
		return github.New(clientKey, secret, callbackURL, scopes...)
	case "apple":
		return apple.New(clientKey, secret, callbackURL, nil, apple.ScopeName, apple.ScopeEmail)
	case "google":
		return google.New(clientKey, secret, callbackURL, scopes...)
	case "yandex":
		return yandex.New(clientKey, secret, callbackURL, scopes...)
	case "azuread":
		return azuread.New(clientKey, secret, callbackURL, nil, scopes...)
	}

	return nil
}
