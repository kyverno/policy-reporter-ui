package oauth

import (
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/amazon"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/gitlab"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
	"golang.org/x/oauth2/yandex"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	oauth2.Config
}

func New(clientID, clientSecret, callbackURL string, endpoint oauth2.Endpoint, scopes []string) (*Authenticator, error) {
	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Endpoint:     endpoint,
		Scopes:       scopes,
	}

	return &Authenticator{
		Config: conf,
	}, nil
}

func (a *Authenticator) GetConfig() oauth2.Config {
	return a.Config
}

func ProviderEndpoint(provider string) (oauth2.Endpoint, error) {
	switch provider {
	case "google":
		return google.Endpoint, nil
	case "github":
		return github.Endpoint, nil
	case "gitlab":
		return gitlab.Endpoint, nil
	case "amazon":
		return amazon.Endpoint, nil
	case "yandex":
		return yandex.Endpoint, nil
	case "microsoft":
		return microsoft.LiveConnectEndpoint, nil
	default:
		return oauth2.Endpoint{}, errors.New("provider not supported, configure the required endpoints manually")
	}
}
