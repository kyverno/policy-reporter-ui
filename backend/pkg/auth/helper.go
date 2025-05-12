package auth

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"go.uber.org/zap"
)

func mapGroups(user goth.User) []string {
	groups := make([]string, 0)

	if GroupClaim == "" {
		return groups
	}

	rawGroups, ok := user.RawData[GroupClaim]
	if !ok {
		keys := make([]string, 0, len(user.RawData))
		for k := range user.RawData {
			keys = append(keys, k)
		}

		zap.L().Warn("group claim not found", zap.Strings("available keys", keys))
		return groups
	}

	mapped, ok := rawGroups.([]any)
	if !ok {
		zap.L().Error("unexpected group claim value, expected []any", zap.String("type", fmt.Sprintf("%T", mapped)))
		return groups
	}

	for _, group := range mapped {
		if g, ok := group.(string); ok {
			groups = append(groups, g)
		}
	}

	zap.L().Debug("mapped user grous", zap.String("user", user.NickName), zap.String("email", user.Email), zap.Strings("groups", groups))

	return groups
}

func GetSession(ctx *gin.Context) sessions.Session {
	s, ok := ctx.Get(sessions.DefaultKey)
	if !ok {
		zap.L().Debug("session not found", zap.String("name", sessions.DefaultKey))
		return nil
	}

	session, ok := s.(sessions.Session)
	if !ok {
		return nil
	}

	return session
}
