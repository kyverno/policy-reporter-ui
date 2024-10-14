package auth

import (
	"github.com/markbates/goth"
)

func mapGroups(user goth.User, claim string) []string {
	groups := make([]string, 0)

	if claim == "" {
		return groups
	}

	rawGroups, ok := user.RawData[claim]
	if !ok {
		return groups
	}

	mapped, ok := rawGroups.([]any)
	if !ok {
		return groups
	}

	for _, group := range mapped {
		if g, ok := group.(string); ok {
			groups = append(groups, g)
		}
	}

	return groups
}
