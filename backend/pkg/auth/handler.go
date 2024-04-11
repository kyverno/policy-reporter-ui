package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

type Handler struct {
	basePath string
}

func NewHandler(basePath string) *Handler {
	return &Handler{
		basePath: basePath,
	}
}

func (h *Handler) Callback(ctx *gin.Context) {
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		zap.L().Error("failed to complete user", zap.Error(err))
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session := sessions.Default(ctx)
	session.Set("profile", NewProfile(user))

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
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		return
	}

	session := sessions.Default(ctx)
	session.Set("profile", NewProfile(user))

	ctx.Redirect(http.StatusTemporaryRedirect, h.basePath)
}

func (h *Handler) Logout(ctx *gin.Context) {
	gothic.Logout(ctx.Writer, ctx.Request)
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

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

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
