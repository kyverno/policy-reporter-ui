package auth

import (
	"fmt"

	"github.com/markbates/goth"
	"go.uber.org/zap"
)

func mapGroups(user goth.User, claim string) []string {
	groups := make([]string, 0)

	if claim == "" {
		return groups
	}

	rawGroups, ok := user.RawData[claim]
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
