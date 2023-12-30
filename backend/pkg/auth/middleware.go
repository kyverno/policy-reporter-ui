package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	profile := sessions.Default(ctx).Get("profile")
	if profile == nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
	} else {
		ctx.Set("profile", profile)
		ctx.Set("token", sessions.Default(ctx).Get("access_token"))
		ctx.Next()
	}
}
