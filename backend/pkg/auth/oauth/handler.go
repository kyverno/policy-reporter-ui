package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kyverno/policy-reporter-ui/pkg/auth"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
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

	session.Set("token", token)
	session.Set("profile", auth.Profile{})
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
	urlParts := strings.Split(h.authenticator.GetConfig().Endpoint.TokenURL, "/")
	urlParts[len(urlParts)-1] = "revoke"

	logoutURL, err := url.Parse(strings.Join(urlParts, "/"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	token := sessions.Default(ctx).Get("token")
	if token == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	parameters := url.Values{}
	parameters.Add("token", token.(*oauth2.Token).AccessToken)
	parameters.Add("client_id", h.authenticator.GetConfig().ClientID)
	parameters.Add("client_secret", h.authenticator.GetConfig().ClientSecret)

	logoutURL.RawQuery = parameters.Encode()

	client := h.authenticator.Client(ctx, token.(*oauth2.Token))
	resp, err := client.Post(logoutURL.String(), "application/x-www-form-urlencoded", nil)
	if err != nil {
		zap.L().Error("failed to revoke token", zap.Error(err), zap.String("host", logoutURL.Host), zap.String("path", logoutURL.Path))
	} else if resp.StatusCode <= 300 {
		content, _ := io.ReadAll(resp.Body)
		zap.L().Info("revoke respose", zap.String("host", logoutURL.Host), zap.String("path", logoutURL.Path), zap.ByteString("body", content))
	}

	sessions.Default(ctx).Clear()

	ctx.Set("profile", nil)
	ctx.Set("token", nil)

	ctx.Redirect(http.StatusTemporaryRedirect, "/login")
}

func (h *Handler) Profile(ctx *gin.Context) {
	token := sessions.Default(ctx).Get("token")
	if token == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	client := h.authenticator.Client(ctx, token.(*oauth2.Token))
	resp, err := client.Get("https://gitlab.example.com/oauth/userinfo")
	if err != nil {
		zap.L().Error("failed to receive userinfo", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	profile := &auth.Profile{}

	if err := json.NewDecoder(resp.Body).Decode(profile); err != nil {
		zap.L().Error("failed to unmarshal userinfo", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, gin.H{
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
