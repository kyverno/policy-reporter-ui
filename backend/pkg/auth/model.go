package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Profile struct {
	SUB       string `json:"sub"`
	Lastname  string `json:"family_name"`
	Firstname string `json:"given_name"`

	Name  string `json:"name"`
	Login string `json:"login"`
	Email string `json:"email"`
}

func (p Profile) GetName() string {
	if p.Name != "" {
		return p.Name
	}

	if p.Firstname != "" || p.Lastname != "" {
		return strings.TrimSpace(p.Firstname + " " + p.Lastname)
	}

	return p.Login
}

func ProfileFrom(ctx *gin.Context) *Profile {
	profile, ok := ctx.Get("profile")
	if !ok {
		return nil
	}

	if p, ok := profile.(Profile); ok {
		return &p
	}

	return nil
}
