package auth

import (
	"encoding/gob"
	"math"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gsessions "github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

const SessionKey = "auth-session"

func Setup(engine *gin.Engine, provider string, tempDir string) {
	gob.Register(Profile{})
	gob.Register(map[string]any{})

	authStore := gsessions.NewFilesystemStore(tempDir, []byte("auth-store"))
	authStore.MaxLength(math.MaxInt64)
	authStore.Options = &gsessions.Options{
		HttpOnly: true,
		MaxAge:   86400 * 30,
		Path:     "/",
	}

	gothic.Store = authStore

	engine.Use(sessions.Sessions(SessionKey, NewStore(authStore)))

	handler := NewHandler()

	engine.GET("/login", Provider(provider), handler.Login)
	engine.GET("/logout", Provider(provider), handler.Logout)
	engine.GET("/profile", Provider(provider), Auth, handler.Profile)
	engine.GET("/callback", Provider(provider), handler.Callback)

	zap.L().Info("setup oauth routes")
}
