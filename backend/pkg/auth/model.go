package auth

import "github.com/gin-gonic/gin"

type Profile struct {
	ID        string `json:"sub"`
	Lastname  string `json:"family_name"`
	Firstname string `json:"given_name"`
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
