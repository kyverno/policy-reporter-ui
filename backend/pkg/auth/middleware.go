package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func Auth(ctx *gin.Context) {
	t := sessions.Default(ctx).Get("token")
	if t == nil {
		abort(ctx)
		return
	}

	token := t.(*oauth2.Token)
	if !token.Valid() {
		abort(ctx)
		return
	}

	profile := sessions.Default(ctx).Get("profile")
	if profile == nil {
		abort(ctx)
		return
	} else {
		ctx.Set("profile", profile)
		ctx.Set("token", token)
		ctx.Next()
	}
}

func abort(ctx *gin.Context) {
	zap.L().Info("abort request", zap.String("path", ctx.Request.URL.Path))

	sessions.Default(ctx).Clear()

	if ctx.Request.URL.Path == "/" {
		ctx.Redirect(http.StatusSeeOther, "/login")
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}
