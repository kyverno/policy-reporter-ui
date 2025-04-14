package auth

import (
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gsessions "github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const SessionKey = "auth-session"

var GroupClaim string = ""

type SessionStorage struct {
	Storage  string
	TempDir  string
	Addr     string
	Database int
	Username string
	Password string
}

func Setup(engine *gin.Engine, basePath, groupKey, provider string, storage SessionStorage) {
	gob.Register(Profile{})
	gob.Register(map[string]any{})

	GroupClaim = groupKey

	if storage.Storage == "filesystem" {
		authStore := gsessions.NewFilesystemStore(storage.TempDir, []byte("auth-store"))
		authStore.MaxLength(2147483647)
		authStore.Options = &gsessions.Options{
			HttpOnly: true,
			MaxAge:   86400 * 30,
			Path:     "/",
		}

		gothic.Store = authStore
		engine.Use(sessions.Sessions(SessionKey, NewStore(authStore)))
	} else if storage.Storage == "redis" {
		authStore := NewRedisStore(&redis.Options{
			Addr:     storage.Addr,
			Username: storage.Username,
			Password: storage.Password,
			DB:       storage.Database,
		}, []byte("auth-store"))

		authStore.MaxLength(2147483647)
		authStore.Options = &gsessions.Options{
			HttpOnly: true,
			MaxAge:   86400 * 30,
			Path:     "/",
		}

		gothic.Store = authStore
		engine.Use(sessions.Sessions(SessionKey, NewStore(authStore)))
	}

	handler := NewHandler(basePath)

	engine.GET("/login", Provider(provider), handler.Login)
	engine.GET("/logout", Provider(provider), handler.Logout)
	engine.GET("/profile", Provider(provider), Auth(basePath), handler.Profile)
	engine.GET("/callback", Provider(provider), handler.Callback)

	zap.L().Info("setup oauth routes")
}
