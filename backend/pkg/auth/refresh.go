package auth

import (
	"errors"
	"hash/fnv"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// Bounded lock striping serializes refreshes for a session within one process.
// Reloading inside the lock lets waiting requests reuse a rotated token.
var refreshLocks [64]sync.Mutex

// refreshExpiredProfile renews an expired OIDC session before either API
// authentication or UI profile validation can discard it.
func refreshExpiredProfile(ctx *gin.Context, profile *Profile) error {
	if profile.ExpiresAt.IsZero() || !profile.ExpiresAt.Before(time.Now()) {
		return nil
	}

	name, err := gothic.GetProviderName(ctx.Request)
	if err != nil || name != "openid-connect" || profile.RefreshToken == "" {
		return errors.New("session cannot be refreshed")
	}

	session := sessions.Default(ctx)
	if session.ID() == "" {
		return errors.New("session has no persistent identity")
	}
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(session.ID()))
	lock := &refreshLocks[hash.Sum64()%uint64(len(refreshLocks))]
	lock.Lock()
	defer lock.Unlock()

	// Bypass the request-local registry: another request may have refreshed
	// this session after the current request loaded its profile.
	stored, err := gothic.Store.New(ctx.Request, SessionKey)
	if err != nil {
		return errors.New("failed to reload session")
	}
	current, ok := stored.Values["profile"].(Profile)
	if !ok || current.ID != profile.ID {
		return errors.New("stored session identity is invalid")
	}
	if current.ExpiresAt.After(time.Now()) {
		session.Set("profile", current)
		*profile = current
		return nil
	}
	if current.RefreshToken == "" {
		return errors.New("stored session cannot be refreshed")
	}
	*profile = current

	provider, err := goth.GetProvider(name)
	if err != nil || !provider.RefreshTokenAvailable() {
		return errors.New("refresh provider unavailable")
	}

	token, err := provider.RefreshToken(profile.RefreshToken)
	if err != nil || token == nil || token.AccessToken == "" || !token.Expiry.After(time.Now()) {
		return errors.New("session refresh failed")
	}

	refreshed := *profile
	refreshed.AccessToken = token.AccessToken
	refreshed.ExpiresAt = token.Expiry
	if token.RefreshToken != "" {
		refreshed.RefreshToken = token.RefreshToken
	}
	if idToken, ok := token.Extra("id_token").(string); ok && idToken != "" {
		refreshed.IDToken = idToken
	}

	// Keep the provider's identity, claims and expiry validation. A successful
	// token exchange alone must not replace the authenticated subject or groups.
	user, err := provider.FetchUser(ProviderSession(name, &refreshed))
	if err != nil || user.UserID == "" || user.UserID != profile.ID {
		return errors.New("refreshed identity is invalid")
	}

	updated := NewProfile(user)
	updated.AssignGroups(mapGroups(user))
	if !updated.ExpiresAt.After(time.Now()) {
		return errors.New("refreshed session is expired")
	}

	session.Set("profile", updated)
	if err := session.Save(); err != nil {
		return errors.New("failed to persist refreshed session")
	}
	*profile = updated
	return nil
}
