package auth

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kyverno/policy-reporter-ui/pkg/cluster"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

type AccessControl struct {
	Emails []string
	Groups []string
}

type Permissions struct {
	AccessControl AccessControl `json:"-"`
}

func (p Permissions) Allowed(profile *Profile) bool {
	if profile == nil {
		return true
	}

	if len(p.AccessControl.Emails) == 0 && len(p.AccessControl.Groups) == 0 {
		return true
	}

	if len(p.AccessControl.Emails) > 0 && slices.Contains(p.AccessControl.Emails, profile.Email) {
		return true
	}

	if len(p.AccessControl.Groups) > 0 && utils.Some(p.AccessControl.Groups, profile.Groups) {
		return true
	}

	return false
}

func NewPermissions(ac cluster.AccessControl) Permissions {
	return Permissions{
		AccessControl: AccessControl{
			Emails: ac.Emails,
			Groups: ac.Groups,
		},
	}
}

func ClusterPermissions(clusters *cluster.Collection) gin.HandlerFunc {
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
			if cl := clusters.Cluster(cluster); cl != nil {
				if access := NewPermissions(cl.AccessControl).Allowed(profile); !access {
					ctx.AbortWithStatus(http.StatusUnauthorized)
					return
				}
			}
		}

		ctx.Next()
	}
}
