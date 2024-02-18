package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Callback(ctx *gin.Context) {
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session := sessions.Default(ctx)
	session.Set("profile", NewProfile(user))

	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *Handler) Login(ctx *gin.Context) {
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		return
	}

	session := sessions.Default(ctx)
	session.Set("profile", NewProfile(user))

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *Handler) Logout(ctx *gin.Context) {
	gothic.Logout(ctx.Writer, ctx.Request)
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	ctx.Redirect(http.StatusTemporaryRedirect, "/login")
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
