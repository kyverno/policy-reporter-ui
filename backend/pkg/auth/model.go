package auth

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"go.uber.org/zap"
)

type Profile struct {
	ID           string    `json:"id"`
	AvatarURL    string    `json:"avatar"`
	NickName     string    `json:"nickname"`
	Lastname     string    `json:"lastname"`
	Firstname    string    `json:"firstname"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Groups       []string  `json:"groups"`
	RefreshToken string    `json:"-"`
	AccessToken  string    `json:"-"`
	IDToken      string    `json:"-"`
	ExpiresAt    time.Time `json:"-"`
}

func (p *Profile) GetName() string {
	if p.Name != "" {
		return p.Name
	}

	if p.NickName != "" {
		return p.NickName
	}

	return p.Email
}

func (p *Profile) AssignGroups(groups []string) {
	p.Groups = groups
}

func NewProfile(user goth.User) Profile {
	return Profile{
		ID:           user.UserID,
		AvatarURL:    user.AvatarURL,
		NickName:     user.NickName,
		Firstname:    user.FirstName,
		Lastname:     user.LastName,
		Name:         user.Name,
		Email:        user.Email,
		ExpiresAt:    user.ExpiresAt,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		IDToken:      user.IDToken,
	}
}

func Session(ctx *gin.Context) (sessions.Session, bool) {
	val, ok := ctx.Get(sessions.DefaultKey)
	if !ok {
		zap.L().Debug("session not found")
		return nil, false
	}

	session, ok := val.(sessions.Session)
	if !ok {
		zap.L().Debug("session is empty")
		return nil, false
	}

	return session, true
}

func ProfileFrom(ctx *gin.Context) *Profile {
	session, ok := Session(ctx)
	if !ok {
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			zap.L().Debug("unable to get profile from session", zap.Any("error", r))
			ClearCookie(ctx)
		}
	}()

	profile := session.Get("profile")

	if profile == nil {
		return nil
	}

	if p, ok := profile.(Profile); ok {
		return &p
	}

	return nil
}

func ClearCookie(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{Name: SessionKey, MaxAge: -1, HttpOnly: true})
}
