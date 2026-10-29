package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gsessions "github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
)

//nolint:paralleltest // Goth/session globals are shared with the existing authentication tests.
func TestCertificateOrOIDCIdentity(t *testing.T) {
	registerTestProfile.Do(func() { gob.Register(Profile{}) })
	oldStore, oldGroup := gothic.Store, GroupClaim
	oldProviders := goth.GetProviders()
	t.Cleanup(func() {
		gothic.Store, GroupClaim = oldStore, oldGroup
		goth.ClearProviders()
		for _, provider := range oldProviders {
			goth.UseProviders(provider)
		}
	})
	issuer := ""
	providerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"issuer": issuer, "authorization_endpoint": issuer + "/authorize", "token_endpoint": issuer + "/token"})
	}))
	defer providerServer.Close()
	issuer = providerServer.URL
	provider, err := openidConnect.New("client", "secret", "https://ui.example/callback", issuer+"/discovery", providerServer.Client(), false)
	if err != nil {
		t.Fatal(err)
	}
	goth.UseProviders(provider)
	claims, _ := json.Marshal(map[string]any{"sub": "oidc-user", "iss": issuer, "aud": "client", "exp": time.Now().Add(time.Hour).Unix(), "email": "oidc@example.test"})
	oidcProfile := Profile{ID: "oidc-user", Email: "oidc@example.test", Groups: []string{"admins"}, ExpiresAt: time.Now().Add(time.Hour), IDToken: "e30." + base64.RawURLEncoding.EncodeToString(claims) + ".signature"}
	store := gsessions.NewFilesystemStore(t.TempDir(), []byte("test-key"))
	gothic.Store = store
	fixture := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	key := fixture.TLS.Certificates[0].PrivateKey.(crypto.Signer)
	fixture.Close()
	template := &x509.Certificate{SerialNumber: big.NewInt(1), NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}, IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}}
	der, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	pair := tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
	roots := x509.NewCertPool()
	roots.AddCert(certificate)
	hash := sha256.Sum256(certificate.Raw)
	mapped := Profile{ID: "certificate-client", Email: "client@example.test", Groups: []string{"readers"}}
	engine := gin.New()
	engine.Use(CertificateIdentities(map[string]Profile{hex.EncodeToString(hash[:]): mapped}))
	engine.Use(sessions.Sessions(SessionKey, NewStore(store)), Provider("openid-connect"))
	engine.GET("/seed", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("profile", oidcProfile)
		if err := session.Save(); err != nil {
			t.Error(err)
		}
	})
	respond := func(c *gin.Context) { c.String(http.StatusOK, ProfileFrom(c).ID) }
	engine.GET("/api", Auth("/"), respond)
	engine.GET("/profile", Auth("/"), NewHandler("/").Profile)
	engine.GET("/ui", Valid("/"), respond)
	engine.GET("/restricted", Auth("/"), func(c *gin.Context) {
		if !(Permissions{AccessControl: AccessControl{Groups: []string{"admins"}}}).Allowed(ProfileFrom(c)) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		respond(c)
	})
	request := func(path string, state *tls.ConnectionState, cookies []*http.Cookie) *httptest.ResponseRecorder {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "https://ui.example"+path, nil)
		req.TLS = state
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, req)
		return recorder
	}
	seed := request("/seed", nil, nil)
	cookie := seed.Result().Cookies()
	verified := &tls.ConnectionState{PeerCertificates: []*x509.Certificate{certificate}, VerifiedChains: [][]*x509.Certificate{{certificate}}}
	for _, path := range []string{"/api", "/ui"} {
		if response := request(path, nil, cookie); response.Code != http.StatusOK || response.Body.String() != "oidc-user" {
			t.Fatalf("OIDC fallback %s: %d %s", path, response.Code, response.Body.String())
		}
		if response := request(path, verified, nil); response.Code != http.StatusOK || response.Body.String() != mapped.ID || len(response.Result().Cookies()) != 0 {
			t.Fatalf("certificate %s: %d %s", path, response.Code, response.Body.String())
		}
		if response := request(path, nil, nil); response.Code == http.StatusOK {
			t.Fatal("certificate authentication persisted without certificate")
		}
	}
	if response := request("/profile", verified, nil); response.Code != http.StatusOK || !strings.Contains(response.Body.String(), mapped.ID) {
		t.Fatalf("certificate frontend profile unavailable: %d %s", response.Code, response.Body.String())
	}
	if response := request("/restricted", verified, cookie); response.Code != http.StatusForbidden {
		t.Fatal("OIDC admin session overrode restricted certificate identity")
	}
	unknown := &x509.Certificate{Raw: []byte("unknown"), NotBefore: certificate.NotBefore, NotAfter: certificate.NotAfter}
	for _, state := range []*tls.ConnectionState{
		{PeerCertificates: []*x509.Certificate{unknown}, VerifiedChains: [][]*x509.Certificate{{unknown}}},
		{PeerCertificates: []*x509.Certificate{certificate}},
	} {
		if response := request("/api", state, cookie); response.Code != http.StatusForbidden {
			t.Fatal("unmapped/unverified certificate fell back to OIDC")
		}
	}
	// The same HTTPS listener accepts the OIDC session or a mapped certificate.
	listener := httptest.NewUnstartedServer(engine)
	listener.TLS = &tls.Config{Certificates: []tls.Certificate{pair}, ClientCAs: roots, ClientAuth: tls.VerifyClientCertIfGiven, MinVersion: tls.VersionTLS12}
	listener.StartTLS()
	defer listener.Close()
	for _, withCertificate := range []bool{false, true} {
		cfg := &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
		if withCertificate {
			cfg.Certificates = []tls.Certificate{pair}
		}
		transport := &http.Transport{TLSClientConfig: cfg}
		defer transport.CloseIdleConnections()
		client := &http.Client{Transport: transport, Timeout: 3 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, listener.URL+"/api", nil)
		if err != nil {
			t.Fatal(err)
		}
		if !withCertificate {
			for _, c := range cookie {
				req.AddCookie(c)
			}
		}
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		expected := "oidc-user"
		if withCertificate {
			expected = mapped.ID
		}
		if response.StatusCode != http.StatusOK || string(body) != expected {
			t.Fatalf("combined listener cert=%v: %d %s", withCertificate, response.StatusCode, body)
		}
	}
	expiredIssuer := &x509.Certificate{NotBefore: time.Now().Add(-2 * time.Hour), NotAfter: time.Now().Add(-time.Hour)}
	chainExpired := &tls.ConnectionState{PeerCertificates: []*x509.Certificate{certificate}, VerifiedChains: [][]*x509.Certificate{{certificate, expiredIssuer}}}
	if response := request("/api", chainExpired, cookie); response.Code != http.StatusForbidden {
		t.Fatal("expired issuer accepted on existing connection")
	}
	certificate.NotAfter = time.Now().Add(-time.Second)
	if response := request("/api", verified, cookie); response.Code != http.StatusForbidden {
		t.Fatal("expired certificate accepted on existing connection")
	}
}
