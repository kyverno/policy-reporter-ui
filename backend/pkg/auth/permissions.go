package auth

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccessControl struct {
	Emails []string
}

type Permissions struct {
	AccessControl AccessControl `json:"-"`
}

func (p Permissions) AllowedEmail(email string) bool {
	if len(p.AccessControl.Emails) == 0 {
		return true
	}

	return slices.Contains(p.AccessControl.Emails, email)
}

func ClusterPermissions(permissions map[string]Permissions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cluster := ctx.Param("cluster")
		if cluster == "" && strings.HasPrefix(ctx.Request.URL.Path, "/proxy/") {
			parts := strings.Split(ctx.Request.URL.Path, "/")
			cluster = strings.TrimSpace(parts[2])
		}

		if cluster == "" {
			ctx.Next()
			return
		}

		if profile := ProfileFrom(ctx); profile != nil {
			if !permissions[cluster].AllowedEmail(profile.Email) {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		ctx.Next()
	}
}
