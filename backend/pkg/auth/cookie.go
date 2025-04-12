package auth

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func NewCookieStore(keyPairs ...[]byte) *CookieStore {
	codes := securecookie.CodecsFromPairs(keyPairs...)
	for _, c := range codes {
		if sc, ok := c.(*securecookie.SecureCookie); ok {
			sc.MaxLength(10000)
		}
	}

	cs := &CookieStore{
		Codecs: codes,
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 30,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		},
	}

	cs.MaxAge(cs.Options.MaxAge)
	return cs
}

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options // default configuration
}

// Get returns a session for the given name after adding it to the registry.
//
// It returns a new session if the sessions doesn't exist. Access IsNew on
// the session to check if it is an existing session or a new one.
//
// It returns a new session and an error if the session exists but could
// not be decoded.
func (s *CookieStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// The difference between New() and Get() is that calling New() twice will
// decode the session data twice, while Get() registers and reuses the same
// decoded session after the first call.
func (s *CookieStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.Values,
			s.Codecs...)
		if err == nil {
			session.IsNew = false
		}
	}
	return session, err
}

// Save adds a single session to the response.
func (s *CookieStore) Save(r *http.Request, w http.ResponseWriter,
	session *sessions.Session,
) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		s.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (s *CookieStore) MaxAge(age int) {
	s.Options.MaxAge = age

	// Set the maxAge for each securecookie instance.
	for _, codec := range s.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}
