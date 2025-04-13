package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

type Handler struct {
	basePath   string
	groupClaim string
}

func NewHandler(basePath, groupClaim string) *Handler {
	return &Handler{basePath: basePath, groupClaim: groupClaim}
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
	profile.AssignGroups(mapGroups(user, h.groupClaim))

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
		ClearCookie(ctx)
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		return
	}

	profile := NewProfile(user)
	profile.AssignGroups(mapGroups(user, h.groupClaim))

	session := sessions.Default(ctx)
	session.Set("profile", profile)

	ctx.Redirect(http.StatusTemporaryRedirect, h.basePath)
}

func (h *Handler) Logout(ctx *gin.Context) {
	if err := gothic.Logout(ctx.Writer, ctx.Request); err != nil {
		zap.L().Error("failed to logout from provider", zap.Error(err))
	}

	session := sessions.Default(ctx)
	session.Clear()
	ClearCookie(ctx)

	if err := session.Save(); err != nil {
		zap.L().Error("failed to save session", zap.Error(err))
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
