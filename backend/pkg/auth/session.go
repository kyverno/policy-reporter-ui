package auth

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/yandex"
)

func GithubSession(profile *Profile) *github.Session {
	return &github.Session{
		AccessToken: profile.AccessToken,
	}
}

func GitlabSession(profile *Profile) *gitlab.Session {
	return &gitlab.Session{
		AccessToken:  profile.AccessToken,
		RefreshToken: profile.RefreshToken,
		ExpiresAt:    profile.ExpiresAt,
	}
}

func AppleSession(profile *Profile) *apple.Session {
	return &apple.Session{
		AccessToken:  profile.AccessToken,
		RefreshToken: profile.RefreshToken,
		ExpiresAt:    profile.ExpiresAt,
	}
}

func GoogleSession(profile *Profile) *google.Session {
	return &google.Session{
		AccessToken:  profile.AccessToken,
		RefreshToken: profile.RefreshToken,
		ExpiresAt:    profile.ExpiresAt,
		IDToken:      profile.IDToken,
	}
}

func AzureadSession(profile *Profile) *azuread.Session {
	return &azuread.Session{
		AccessToken:  profile.AccessToken,
		RefreshToken: profile.RefreshToken,
		ExpiresAt:    profile.ExpiresAt,
	}
}

func YandexSession(profile *Profile) *yandex.Session {
	return &yandex.Session{
		AccessToken:  profile.AccessToken,
		RefreshToken: profile.RefreshToken,
		ExpiresAt:    profile.ExpiresAt,
	}
}

func OIDCSession(profile *Profile) *openidConnect.Session {
	return &openidConnect.Session{
		AccessToken:  profile.AccessToken,
		RefreshToken: profile.RefreshToken,
		ExpiresAt:    profile.ExpiresAt,
		IDToken:      profile.IDToken,
	}
}

func AmazonSession(profile *Profile) *amazon.Session {
	return &amazon.Session{
		AccessToken:  profile.AccessToken,
		RefreshToken: profile.RefreshToken,
		ExpiresAt:    profile.ExpiresAt,
	}
}

func ProviderSession(provider string, profile *Profile) goth.Session {
	switch provider {
	case "amazon":
		return AmazonSession(profile)
	case "gitlab":
		return GitlabSession(profile)
	case "github":
		return GithubSession(profile)
	case "apple":
		return AppleSession(profile)
	case "google":
		return GoogleSession(profile)
	case "yandex":
		return YandexSession(profile)
	case "azuread":
		return AzureadSession(profile)
	case "openid-connect":
		return OIDCSession(profile)
	}

	return nil
}
