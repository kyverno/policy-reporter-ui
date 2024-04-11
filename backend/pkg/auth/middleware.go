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

func Valid(ctx *gin.Context) {
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
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	session := sessions.Default(ctx)

	sess := ProviderSession(providerName, profile)
	if sess == nil {
		zap.L().Error("could not create session from profile", zap.Error(err))

		logout(ctx)
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	user, err := provider.FetchUser(sess)
	if err != nil {
		zap.L().Error("failed to validate session", zap.Error(err))

		logout(ctx)
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	session.Set("profile", NewProfile(user))
	session.Save()

	ctx.Next()
}

func Auth(ctx *gin.Context) {
	profile := ProfileFrom(ctx)

	if profile == nil {
		abort(ctx, "")
		return
	}

	if !profile.ExpiresAt.IsZero() && profile.ExpiresAt.Before(time.Now()) {
		abort(ctx, "expired session")
		return
	}

	ctx.Next()
}

func abort(ctx *gin.Context, err string) {
	if err != "" {
		zap.L().Info("abort request", zap.String("path", ctx.Request.URL.Path), zap.String("err", err))
	}

	logout(ctx)

	if ctx.Request.URL.Path == "/" {
		ctx.Redirect(http.StatusSeeOther, "login")
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}

func logout(ctx *gin.Context) error {
	gothic.Logout(ctx.Writer, ctx.Request)
	session := sessions.Default(ctx)
	session.Clear()

	return session.Save()
}
