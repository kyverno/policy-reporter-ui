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

var (
	GroupClaim string = ""
	keyPair           = []byte("auth-store")
)

type SessionStorage struct {
	Storage  string
	TempDir  string
	Addr     string
	Database int
	Username string
	Password string
}

func Setup(engine *gin.Engine, basePath, groupKey, provider string, s SessionStorage) {
	gob.Register(Profile{})
	gob.Register(map[string]any{})

	GroupClaim = groupKey

	switch s.Storage {
	case "redis":
		authStore := NewRedisStore(&redis.Options{
			Addr:     s.Addr,
			Username: s.Username,
			Password: s.Password,
			DB:       s.Database,
		}, keyPair)

		authStore.MaxLength(2147483647)
		authStore.Options = &gsessions.Options{
			HttpOnly: true,
			MaxAge:   86400 * 30,
			Path:     "/",
		}

		gothic.Store = authStore
		engine.Use(sessions.Sessions(SessionKey, NewStore(authStore)))
	default:
		authStore := gsessions.NewFilesystemStore(s.TempDir, keyPair)
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
