package auth

import (
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gsessions "github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

const SessionKey = "auth-session"

type SessionStorage struct {
	Type    string
	TempDir string
}

func Setup(engine *gin.Engine, basePath, groupKey, provider string, storage SessionStorage) {
	gob.Register(Profile{})
	gob.Register(map[string]any{})

	if storage.Type == "filesystem" {
		authStore := gsessions.NewFilesystemStore(storage.TempDir, []byte("auth-store"))
		authStore.MaxLength(2147483647)
		authStore.Options = &gsessions.Options{
			HttpOnly: true,
			MaxAge:   86400 * 30,
			Path:     "/",
		}

		gothic.Store = authStore
		engine.Use(sessions.Sessions(SessionKey, NewStore(authStore)))
	} else {
		authStore := gsessions.NewCookieStore([]byte("auth-store"))
		authStore.Options = &gsessions.Options{
			HttpOnly: true,
			MaxAge:   86400 * 30,
			Path:     "/",
		}

		gothic.Store = authStore
		engine.Use(sessions.Sessions(SessionKey, NewStore(authStore)))
	}

	handler := NewHandler(basePath, groupKey)

	engine.GET("/login", Provider(provider), handler.Login)
	engine.GET("/logout", Provider(provider), handler.Logout)
	engine.GET("/profile", Provider(provider), Auth(basePath), handler.Profile)
	engine.GET("/callback", Provider(provider), handler.Callback)

	zap.L().Info("setup oauth routes")
}
