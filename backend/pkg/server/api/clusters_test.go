package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/kyverno/policy-reporter-ui/pkg/auth"
	"github.com/kyverno/policy-reporter-ui/pkg/cluster"
	"github.com/kyverno/policy-reporter-ui/pkg/customboard"
)

// Use the real Core client against small HTTP responses, with no client mocks.
func clusterHandler(t *testing.T, boards ...*customboard.CustomBoard) (*Handler, *atomic.Int32) {
	t.Helper()
	calls := new(atomic.Int32)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v2/sources":
			_, _ = w.Write([]byte(`["kyverno"]`))
		case "/v2/sources/categories", "/v2/resource/resource/status-counts", "/v2/resource/resource/severity-counts", "/v2/resource/resource/source-categories":
			_, _ = w.Write([]byte(`[]`))
		case "/v2/targets", "/v2/resource/resource", "/v2/namespace-scoped/results":
			_, _ = w.Write([]byte(`{}`))
		default:
			t.Errorf("unexpected Core request: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(upstream.Close)
	clients := cluster.NewCollection(cluster.Config{Name: "Production", Host: upstream.URL}, cluster.Config{Name: "Staging", Host: upstream.URL})
	return NewHandler(&Config{}, clients, customboard.NewCollection(boards...)), calls
}

func boardRequest(handler gin.HandlerFunc, clusterID, id string, profile *auth.Profile) *httptest.ResponseRecorder {
	router := gin.New()
	if profile != nil {
		router.Use(sessions.Sessions("test", cookie.NewStore([]byte("test-session-secret"))))
		router.Use(func(ctx *gin.Context) { sessions.Default(ctx).Set("profile", *profile); ctx.Next() })
	}
	router.GET("/:cluster/:id", func(ctx *gin.Context) {
		ctx.Params = append(ctx.Params, gin.Param{Key: "resource", Value: "resource"})
		handler(ctx)
	})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/"+clusterID+"/"+id, nil))
	return response
}

func TestCustomBoardClusterMismatch(t *testing.T) {
	t.Parallel()
	// These are all seven handlers registered under /:cluster/custom-board.
	for _, tt := range []struct {
		name    string
		handler func(*Handler, *gin.Context)
	}{
		{"dashboard", (*Handler).GetCustomBoard},
		{"cluster-resource-results", (*Handler).ListCustomBoardClusterResourceResults},
		{"resource-results", (*Handler).ListCustomBoardResourceResults},
		{"cluster-results", (*Handler).ListCustomBoardClusterScopedResults},
		{"results", (*Handler).ListCustomBoardNamespaceScopedResults},
		{"resource-details", (*Handler).GetCustomBoardResourceDetails},
		{"resource-detailed-results", (*Handler).ListCustomBoardResourceDetailedResults},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h, calls := clusterHandler(t, &customboard.CustomBoard{ID: "board", Clusters: []string{"staging"}})
			response := boardRequest(func(ctx *gin.Context) { tt.handler(h, ctx) }, "production", "board", nil)
			if response.Code != http.StatusNotFound {
				t.Errorf("status = %d, want 404", response.Code)
			}
			if got := calls.Load(); got != 0 {
				t.Errorf("Core requests = %d, want 0", got)
			}
		})
	}
}

func TestCustomBoardApplicable(t *testing.T) {
	t.Parallel()
	for _, clusters := range [][]string{nil, {}, {"production"}} {
		h, calls := clusterHandler(t, &customboard.CustomBoard{ID: "board", Name: "Board", Clusters: clusters})
		for _, handler := range []gin.HandlerFunc{h.GetCustomBoard, h.ListCustomBoardNamespaceScopedResults, h.GetCustomBoardResourceDetails} {
			calls.Store(0)
			response := boardRequest(handler, "production", "board", nil)
			if response.Code != http.StatusOK {
				t.Fatalf("clusters %v: status = %d, body = %s", clusters, response.Code, response.Body)
			}
			if calls.Load() == 0 {
				t.Fatal("expected a Core request")
			}
		}
	}
}

func TestCustomBoardMissingAndDenied(t *testing.T) {
	t.Parallel()
	h, calls := clusterHandler(t, &customboard.CustomBoard{ID: "board", Clusters: []string{"production"}, AccessControl: customboard.AccessControl{Emails: []string{"allowed@example.org"}}})
	for _, tt := range []struct {
		name, cluster, board string
		profile              *auth.Profile
		status               int
	}{
		{"missing-board", "production", "missing", nil, http.StatusNotFound},
		{"missing-cluster", "missing", "board", nil, http.StatusNotFound},
		{"denied", "production", "board", &auth.Profile{Email: "denied@example.org"}, http.StatusUnauthorized},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			response := boardRequest(h.GetCustomBoard, tt.cluster, tt.board, tt.profile)
			if response.Code != tt.status {
				t.Errorf("status = %d, want %d", response.Code, tt.status)
			}
			if calls.Load() != 0 {
				t.Fatal("unexpected Core request")
			}
		})
	}
}

func TestLayoutClusterVisibility(t *testing.T) {
	t.Parallel()
	h, _ := clusterHandler(t,
		&customboard.CustomBoard{ID: "restricted", Name: "A Restricted", Clusters: []string{"production"}},
		&customboard.CustomBoard{ID: "unrestricted", Name: "Z Unrestricted"},
		&customboard.CustomBoard{ID: "denied", Name: "B Denied", Clusters: []string{"production"}, AccessControl: customboard.AccessControl{Emails: []string{"allowed@example.org"}}},
	)
	for _, tt := range []struct {
		name, cluster string
		profile       *auth.Profile
		titles        []string
	}{
		{"anonymous-matching", "production", nil, []string{"A Restricted", "B Denied", "Z Unrestricted"}},
		{"anonymous-nonmatching", "staging", nil, []string{"Z Unrestricted"}},
		{"authenticated-matching", "production", &auth.Profile{Email: "denied@example.org"}, []string{"A Restricted", "Z Unrestricted"}},
		{"authenticated-nonmatching", "staging", &auth.Profile{Email: "allowed@example.org"}, []string{"Z Unrestricted"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			response := boardRequest(h.Layout, tt.cluster, "unused", tt.profile)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d: %s", response.Code, response.Body)
			}
			var layout struct {
				CustomBoards []struct {
					Title string `json:"title"`
				} `json:"customBoards"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &layout); err != nil {
				t.Fatal(err)
			}
			titles := make([]string, 0, len(layout.CustomBoards))
			for _, board := range layout.CustomBoards {
				titles = append(titles, board.Title)
			}
			if !reflect.DeepEqual(titles, tt.titles) {
				t.Errorf("titles = %v, want %v", titles, tt.titles)
			}
		})
	}
}
