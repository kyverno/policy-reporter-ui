package auth

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
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
		return nil, false
	}

	session, ok := val.(sessions.Session)
	if !ok {
		return nil, false
	}

	return session, true
}

func ProfileFrom(ctx *gin.Context) *Profile {
	session, ok := Session(ctx)
	if !ok {
		return nil
	}
	defer func() *Profile {
		recover()
		return nil
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
