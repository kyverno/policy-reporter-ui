package auth

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gsessions "github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
)

var registerTestProfile sync.Once

type failingRefreshStore struct {
	gsessions.Store
	fail bool
}

func (s *failingRefreshStore) Get(r *http.Request, name string) (*gsessions.Session, error) {
	original, err := s.Store.Get(r, name)
	if err != nil {
		return nil, err
	}
	session := gsessions.NewSession(s, name)
	session.Values, session.Options = original.Values, original.Options
	session.ID, session.IsNew = original.ID, original.IsNew
	return session, nil
}
func (s *failingRefreshStore) Save(r *http.Request, w http.ResponseWriter, session *gsessions.Session) error {
	if s.fail {
		return errors.New("save failed")
	}
	return s.Store.Save(r, w, session)
}

func TestOIDCSessionRefresh(t *testing.T) {
	registerTestProfile.Do(func() { gob.Register(Profile{}) })
	gin.SetMode(gin.TestMode)

	for _, tc := range []struct {
		name          string
		path          string
		fresh         bool
		missingToken  bool
		rejectRefresh bool
		rotate        bool
		omitIDToken   bool
		wrongSubject  bool
		invalidIssuer bool
		saveFailure   bool
		concurrent    bool
		validOldID    bool
		wantStatus    int
		wantRefresh   int
	}{
		{name: "concurrent API refresh", path: "/api", rotate: true, concurrent: true, wantStatus: http.StatusOK, wantRefresh: 1},
		{name: "missing new ID token with valid old ID token", path: "/api", omitIDToken: true, validOldID: true, wantStatus: http.StatusOK, wantRefresh: 1},
		{name: "expired API session", path: "/api", rotate: true, wantStatus: http.StatusOK, wantRefresh: 1},
		{name: "expired UI session", path: "/ui", rotate: true, wantStatus: http.StatusOK, wantRefresh: 1},
		{name: "unrotated refresh token", path: "/api", wantStatus: http.StatusOK, wantRefresh: 1},
		{name: "fresh API session", path: "/api", fresh: true, wantStatus: http.StatusOK},
		{name: "no refresh token", path: "/api", missingToken: true, wantStatus: http.StatusUnauthorized},
		{name: "refresh rejected", path: "/api", rejectRefresh: true, wantStatus: http.StatusUnauthorized, wantRefresh: 2},
		{name: "UI refresh rejected", path: "/ui", rejectRefresh: true, wantStatus: http.StatusTemporaryRedirect, wantRefresh: 2},
		{name: "expired ID token without replacement", path: "/api", omitIDToken: true, wantStatus: http.StatusUnauthorized, wantRefresh: 1},
		{name: "subject changed", path: "/api", wrongSubject: true, wantStatus: http.StatusUnauthorized, wantRefresh: 1},
		{name: "invalid issuer", path: "/api", invalidIssuer: true, wantStatus: http.StatusUnauthorized, wantRefresh: 1},
		{name: "session save failed", path: "/api", saveFailure: true, wantStatus: http.StatusUnauthorized, wantRefresh: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oldStore, oldGroupClaim := gothic.Store, GroupClaim
			previous, previousErr := goth.GetProvider("openid-connect")
			t.Cleanup(func() {
				gothic.Store, GroupClaim = oldStore, oldGroupClaim
				goth.ClearProviders()
				if previousErr == nil {
					goth.UseProviders(previous)
				}
			})
			GroupClaim = "groups"

			var refreshes atomic.Int32
			var issuer string
			token := func(subject, tokenIssuer string, expiry time.Time) string {
				claims, err := json.Marshal(map[string]any{
					"sub": subject, "iss": tokenIssuer, "aud": "client", "exp": expiry.Unix(), "groups": []string{"readers"},
				})
				if err != nil {
					t.Fatal(err)
				}
				return "e30." + base64.RawURLEncoding.EncodeToString(claims) + ".signature"
			}
			providerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/discovery":
					_ = json.NewEncoder(w).Encode(map[string]string{"issuer": issuer, "authorization_endpoint": issuer + "/authorize", "token_endpoint": issuer + "/token"})
				case "/token":
					refreshes.Add(1)
					if err := r.ParseForm(); err != nil {
						t.Error(err)
					}
					if r.Form.Get("grant_type") != "refresh_token" || r.Form.Get("refresh_token") != "old-refresh" {
						t.Errorf("unexpected refresh request: %v", r.Form)
					}
					if tc.rejectRefresh {
						w.WriteHeader(http.StatusBadRequest)
						_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
						return
					}
					subject, tokenIssuer := "user", issuer
					if tc.wrongSubject {
						subject = "another-user"
					}
					if tc.invalidIssuer {
						tokenIssuer = "wrong-issuer"
					}
					response := map[string]any{"access_token": "new-access", "token_type": "Bearer", "expires_in": 7200}
					if tc.rotate {
						response["refresh_token"] = "new-refresh"
					}
					if !tc.omitIDToken {
						response["id_token"] = token(subject, tokenIssuer, time.Now().Add(2*time.Minute))
					}
					_ = json.NewEncoder(w).Encode(response)
				default:
					http.NotFound(w, r)
				}
			}))
			defer providerServer.Close()
			issuer = providerServer.URL
			provider, err := openidConnect.New("client", "secret", issuer+"/callback", issuer+"/discovery", providerServer.Client(), false)
			if err != nil {
				t.Fatal(err)
			}
			goth.UseProviders(provider)

			profile := Profile{ID: "user", AccessToken: "old-access", RefreshToken: "old-refresh", ExpiresAt: time.Now().Add(-10 * time.Minute), IDToken: token("user", issuer, time.Now().Add(-10*time.Minute)), Groups: []string{"old-group"}}
			if tc.validOldID {
				profile.IDToken = token("user", issuer, time.Now().Add(time.Hour))
			}
			if tc.fresh {
				profile.ExpiresAt = time.Now().Add(time.Hour)
			}
			if tc.missingToken {
				profile.RefreshToken = ""
			}
			store := &failingRefreshStore{Store: gsessions.NewFilesystemStore(t.TempDir(), []byte("test-session-key"))}
			gothic.Store = store
			router := gin.New()
			router.Use(sessions.Sessions(SessionKey, NewStore(store)), Provider("openid-connect"))
			router.GET("/seed", func(c *gin.Context) {
				s := sessions.Default(c)
				s.Set("profile", profile)
				if err := s.Save(); err != nil {
					t.Error(err)
				}
				c.Status(http.StatusOK)
			})
			respond := func(c *gin.Context) { c.JSON(http.StatusOK, ProfileFrom(c)) }
			if tc.concurrent {
				var arrived sync.WaitGroup
				arrived.Add(2)
				router.Use(func(c *gin.Context) {
					if c.Request.URL.Path == "/api" {
						_ = ProfileFrom(c)
						arrived.Done()
						arrived.Wait()
					}
					c.Next()
				})
			}
			router.GET("/api", Auth("/"), respond)
			router.GET("/ui", Valid("/"), respond)
			router.GET("/stored", respond)
			request := func(path string, cookies []*http.Cookie) *httptest.ResponseRecorder {
				r := httptest.NewRequest(http.MethodGet, path, nil)
				for _, cookie := range cookies {
					r.AddCookie(cookie)
				}
				w := httptest.NewRecorder()
				router.ServeHTTP(w, r)
				return w
			}
			seed := request("/seed", nil)
			store.fail = tc.saveFailure
			var response *httptest.ResponseRecorder
			if tc.concurrent {
				responses := make(chan *httptest.ResponseRecorder, 2)
				for range 2 {
					go func() { responses <- request(tc.path, seed.Result().Cookies()) }()
				}
				response = <-responses
				if other := <-responses; other.Code != http.StatusOK {
					t.Fatalf("concurrent request status = %d", other.Code)
				}
			} else {
				response = request(tc.path, seed.Result().Cookies())
			}
			if response.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, tc.wantStatus, response.Body.String())
			}
			if int(refreshes.Load()) != tc.wantRefresh {
				t.Errorf("refresh requests = %d, want %d", refreshes.Load(), tc.wantRefresh)
			}
			if tc.wantStatus == http.StatusOK && !tc.fresh {
				cookies := response.Result().Cookies()
				if len(cookies) == 0 {
					cookies = seed.Result().Cookies()
				}
				stored := request("/stored", cookies)
				var actual Profile
				if err := json.Unmarshal(stored.Body.Bytes(), &actual); err != nil {
					t.Fatal(err)
				}
				if actual.ID != "user" || len(actual.Groups) != 1 || actual.Groups[0] != "readers" {
					t.Errorf("refreshed profile not persisted: %+v", actual)
				}
				// Read the server-side session to check secrets that JSON deliberately omits.
				check := httptest.NewRequest(http.MethodGet, "/", nil)
				for _, cookie := range cookies {
					check.AddCookie(cookie)
				}
				saved, err := store.Get(check, SessionKey)
				if err != nil {
					t.Fatal(err)
				}
				updated := saved.Values["profile"].(Profile)
				wantRefresh := "old-refresh"
				if tc.rotate {
					wantRefresh = "new-refresh"
				}
				if updated.AccessToken != "new-access" || updated.RefreshToken != wantRefresh || (!tc.omitIDToken && updated.IDToken == profile.IDToken) || !updated.ExpiresAt.After(time.Now()) {
					t.Error("token fields were not refreshed and persisted")
				}
			}
		})
	}
}
