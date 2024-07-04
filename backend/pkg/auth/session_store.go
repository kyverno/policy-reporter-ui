package auth

import (
	"github.com/gin-contrib/sessions"
	gsessions "github.com/gorilla/sessions"
)

func NewStore(s gsessions.Store) sessions.Store {
	return &store{s}
}

type store struct {
	gsessions.Store
}

func (c *store) Options(_ sessions.Options) {}
