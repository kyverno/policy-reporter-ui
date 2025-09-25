package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
	"go.uber.org/zap"
)

type Handler struct {
	basePath string
}

func NewHandler(basePath string) *Handler {
	return &Handler{basePath: basePath}
}

func (h *Handler) Callback(ctx *gin.Context) {
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ClearCookie(ctx)
		zap.L().Error("failed to complete user", zap.Error(err))
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	profile := NewProfile(user)
	profile.AssignGroups(mapGroups(user))

	session := sessions.Default(ctx)
	session.Set("profile", profile)

	if err := session.Save(); err != nil {
		zap.L().Error("failed to save session", zap.Error(err))
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, h.basePath)
}

func (h *Handler) Login(ctx *gin.Context) {
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		zap.L().Debug("failed to complete user auth, restart auth process", zap.Error(err))
		ClearCookie(ctx)
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		return
	}

	zap.L().Debug("auth completed", zap.String("user", user.UserID))
	profile := NewProfile(user)
	profile.AssignGroups(mapGroups(user))

	session := sessions.Default(ctx)
	session.Set("profile", profile)

	ctx.Redirect(http.StatusTemporaryRedirect, h.basePath)
}

func (h *Handler) Logout(ctx *gin.Context) {
	endSessionEndopoint := GetEndSessionEndpoint(ctx.Request)

	if err := gothic.Logout(ctx.Writer, ctx.Request); err != nil {
		zap.L().Error("failed to logout from provider", zap.Error(err))
	}

	session := sessions.Default(ctx)
	session.Clear()

	if err := session.Save(); err != nil {
		zap.L().Error("failed to save session", zap.Error(err))
	}

	ClearCookie(ctx)

	if endSessionEndopoint != "" {
		ctx.Redirect(http.StatusTemporaryRedirect, endSessionEndopoint)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, h.basePath+"login")
}

func (h *Handler) Profile(ctx *gin.Context) {
	profile := ProfileFrom(ctx)
	if profile == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":   profile.ID,
		"name": profile.GetName(),
	})
}

func GetEndSessionEndpoint(req *http.Request) string {
	providerName, err := gothic.GetProviderName(req)
	if err != nil {
		return ""
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return ""
	}

	if p, ok := provider.(*openidConnect.Provider); ok {
		return p.OpenIDConfig.EndSessionEndpoint
	}

	return ""
}
