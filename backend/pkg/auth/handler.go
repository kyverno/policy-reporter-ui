package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	authenticator *Authenticator
}

func NewHandler(auth *Authenticator) *Handler {
	return &Handler{authenticator: auth}
}

func (h *Handler) Callback(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if ctx.Query("state") != session.Get("state") {
		ctx.String(http.StatusBadRequest, "Invalid state parameter.")
		return
	}

	// Exchange an authorization code for a token.
	token, err := h.authenticator.Exchange(ctx.Request.Context(), ctx.Query("code"))
	if err != nil {
		zap.L().Error("failed to exchange an authorization code", zap.String("code", ctx.Query("code")), zap.Error(err))

		ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
		return
	}

	idToken, err := h.authenticator.VerifyIDToken(ctx.Request.Context(), token)
	if err != nil {
		zap.L().Error("failed to verify ID Token", zap.Error(err))

		ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
		return
	}

	var profile Profile
	if err := idToken.Claims(&profile); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)
	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *Handler) Login(ctx *gin.Context) {
	state, err := generateRandomState()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Save the state inside the session.
	session := sessions.Default(ctx)
	session.Set("state", state)
	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, h.authenticator.AuthCodeURL(state))
}

func (h *Handler) Logout(ctx *gin.Context) {
	logoutURL, err := url.Parse(strings.TrimSuffix(h.authenticator.Config.Endpoint.AuthURL, "/auth") + "/logout")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", h.authenticator.ClientID)
	logoutURL.RawQuery = parameters.Encode()

	ctx.Redirect(http.StatusTemporaryRedirect, logoutURL.String())
}

func (h *Handler) Profile(ctx *gin.Context) {
	profile := ProfileFrom(ctx)
	if profile == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":        profile.ID,
		"firstname": profile.Firstname,
		"lastname":  profile.Lastname,
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
