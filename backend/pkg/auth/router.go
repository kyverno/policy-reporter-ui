package auth

import (
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const SessionKey = "auth-session"

func Setup(engine *gin.Engine, authenticator *Authenticator) {
	gob.Register(Profile{})

	store := cookie.NewStore([]byte("secret"))
	engine.Use(sessions.Sessions(SessionKey, store))

	handler := NewHandler(authenticator)

	engine.GET("/login", handler.Login)
	engine.GET("/logout", Auth, handler.Logout)
	engine.GET("/profile", Auth, handler.Profile)
	engine.GET("/callback", handler.Callback)

	zap.L().Info("setup oauth routes")
}
