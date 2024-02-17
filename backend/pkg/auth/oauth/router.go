package oauth

import (
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/kyverno/policy-reporter-ui/pkg/auth"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const SessionKey = "auth-session"

func Setup(engine *gin.Engine, authenticator *Authenticator) {
	gob.Register(auth.Profile{})
	gob.Register(&oauth2.Token{})

	store := cookie.NewStore([]byte("secret"))
	engine.Use(sessions.Sessions(SessionKey, store))

	handler := NewHandler(authenticator)

	engine.GET("/login", handler.Login)
	engine.GET("/logout", auth.Auth, handler.Logout)
	engine.GET("/profile", auth.Auth, handler.Profile)
	engine.GET("/callback", handler.Callback)

	zap.L().Info("setup oauth routes")
}
