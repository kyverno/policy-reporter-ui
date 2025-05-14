package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

const (
	ProviderKey = "provider"
)

func Provider(provider string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//nolint:staticcheck
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx, ProviderKey, provider))
		ctx.Set(ProviderKey, provider)
		ctx.Next()
	}
}

func Valid(basePath string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providerName, err := gothic.GetProviderName(ctx.Request)
		if err != nil {
			zap.L().Error("failed to get provider name", zap.Error(err))
			_ = ctx.AbortWithError(http.StatusPreconditionFailed, errors.New("provider name not avaialable in request"))
			return
		}

		provider, err := goth.GetProvider(providerName)
		if err != nil {
			zap.L().Error("failed to get requested provider", zap.Error(err))
			_ = ctx.AbortWithError(http.StatusPreconditionFailed, errors.New("provider not available"))
			return
		}

		profile := ProfileFrom(ctx)
		if profile == nil {
			zap.L().Error("profile not found", zap.Error(err))

			logout(ctx)
			ctx.Redirect(http.StatusTemporaryRedirect, basePath+"login")
			return
		}

		session := sessions.Default(ctx)

		sess := ProviderSession(providerName, profile)
		if sess == nil {
			zap.L().Error("could not create session from profile", zap.Error(err))

			logout(ctx)
			ctx.Redirect(http.StatusTemporaryRedirect, basePath+"login")
			return
		}

		user, err := provider.FetchUser(sess)
		if err != nil {
			zap.L().Error("failed to validate session", zap.Error(err))

			logout(ctx)
			ctx.Redirect(http.StatusTemporaryRedirect, basePath+"login")
			return
		}

		newProfile := NewProfile(user)
		newProfile.AssignGroups(mapGroups(user))

		session.Set("profile", newProfile)
		if err := session.Save(); err != nil {
			zap.L().Error("failed to save profile session", zap.Error(err))
		}

		ctx.Next()
	}
}

func Auth(basePath string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		profile := ProfileFrom(ctx)

		if profile == nil {
			abort(ctx, basePath, "")
			return
		}

		if !profile.ExpiresAt.IsZero() && profile.ExpiresAt.Before(time.Now()) {
			abort(ctx, basePath, "expired session")
			return
		}

		ctx.Next()
	}
}

func abort(ctx *gin.Context, basePath, err string) {
	if err != "" {
		zap.L().Info("abort request", zap.String("path", ctx.Request.URL.Path), zap.String("err", err))
	}

	logout(ctx)

	if ctx.Request.URL.Path == "/" {
		zap.L().Debug("redirect request", zap.String("path", ctx.Request.URL.Path), zap.String("basePath", basePath), zap.String("target", basePath+"login"))
		ctx.Redirect(http.StatusSeeOther, basePath+"login")
	} else {
		zap.L().Debug("abort request", zap.String("path", ctx.Request.URL.Path), zap.String("basePath", basePath))
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}

func logout(ctx *gin.Context) {
	if err := gothic.Logout(ctx.Writer, ctx.Request); err != nil {
		zap.L().Error("failed to logout from provider", zap.Error(err))
		ClearCookie(ctx)
	}

	session := GetSession(ctx)
	if session == nil {
		return
	}

	session.Clear()

	if err := session.Save(); err != nil {
		zap.L().Error("failed to save session", zap.Error(err))
	}
}
