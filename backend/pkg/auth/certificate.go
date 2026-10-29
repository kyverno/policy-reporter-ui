package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const certificateProfileKey = "verified-certificate-profile"

func certificateProfile(ctx *gin.Context) *Profile {
	value, _ := ctx.Get(certificateProfileKey)
	profile, _ := value.(*Profile)
	return profile
}

// CertificateIdentities authenticates only verified, explicitly mapped certificates.
func CertificateIdentities(identities map[string]Profile) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		connection := ctx.Request.TLS
		if connection == nil || len(connection.PeerCertificates) == 0 {
			ctx.Next()
			return
		}
		if len(connection.VerifiedChains) == 0 || len(connection.VerifiedChains[0]) == 0 {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		chain := connection.VerifiedChains[0]
		// Connections can outlive certificate validity; never carry this identity in a session cookie.
		now := time.Now()
		for _, certificate := range chain {
			if now.Before(certificate.NotBefore) || now.After(certificate.NotAfter) {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
		fingerprint := sha256.Sum256(chain[0].Raw)
		profile, ok := identities[hex.EncodeToString(fingerprint[:])]
		if !ok {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ctx.Set(certificateProfileKey, &profile)
		ctx.Next()
	}
}
