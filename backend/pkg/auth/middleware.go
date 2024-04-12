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

func Provider(provider string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx, "provider", provider))
		ctx.Set("provider", provider)
		ctx.Next()
	}
}

func Valid(basePath string) gin.HandlerFunc {
	zap.L().Debug("URL base path", zap.String("basePath", basePath))

	return func(ctx *gin.Context) {
		providerName, err := gothic.GetProviderName(ctx.Request)
		if err != nil {
			zap.L().Error("failed to get provider name", zap.Error(err))
			ctx.AbortWithError(http.StatusPreconditionFailed, errors.New("provider name not avaialable in request"))
			return
		}

		provider, err := goth.GetProvider(providerName)
		if err != nil {
			zap.L().Error("failed to get requested provider", zap.Error(err))
			ctx.AbortWithError(http.StatusPreconditionFailed, errors.New("provider not available"))
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

		session.Set("profile", NewProfile(user))
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
		zap.L().Debug("request URL", zap.Reflect("URL", ctx.Request.URL))
		ctx.Redirect(http.StatusSeeOther, basePath+"login")
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}

func logout(ctx *gin.Context) error {
	gothic.Logout(ctx.Writer, ctx.Request)
	session := sessions.Default(ctx)
	session.Clear()

	err := session.Save()
	if err != nil {
		zap.L().Error("failed to save session", zap.Error(err))
	}

	return err
}
