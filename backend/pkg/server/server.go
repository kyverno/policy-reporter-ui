package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.uber.org/zap"

	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/model"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/customboard"
	"github.com/kyverno/policy-reporter-ui/pkg/server/api"
)

type APIHandler interface {
	Register(*gin.RouterGroup)
}

type Server struct {
	middelware []gin.HandlerFunc
	apis       map[string]*model.Endpoints
	engine     *gin.Engine
	api        *gin.RouterGroup
	proxies    *gin.RouterGroup
	port       int
}

func (s *Server) Start() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.port))
}

func (s *Server) RegisterUI(path string, middleware []gin.HandlerFunc) {
	fileServer := http.FileServer(http.Dir(path))

	handler := append(s.middelware, middleware...)

	s.engine.NoRoute(append(handler, func(c *gin.Context) {
		zap.L().Debug("serving static file server", zap.String("path", c.Request.URL.Path), zap.Any("url", c.Request.URL))
		fileServer.ServeHTTP(c.Writer, c.Request)
	})...)
}

func (s *Server) RegisterCluster(name string, client *core.Client, plugins map[string]*plugin.Client, proxy *httputil.ReverseProxy) {
	id := slug.Make(name)

	s.apis[id] = &model.Endpoints{Name: name, Core: client, Plugins: plugins}
	group := s.proxies.Group(id)

	group.Group("core").Any("/*proxy", func(ctx *gin.Context) {
		req := ctx.Request.Clone(ctx)
		req.URL.Path = strings.TrimPrefix(ctx.Param("proxy"), "/")

		proxy.ServeHTTP(ctx.Writer, req)
	})

	zap.L().Debug("cluster registered", zap.String("name", name), zap.String("id", id))
}

func (s *Server) RegisterAPI(c *api.Config, customBoards *customboard.Collection) {
	handler := api.NewHandler(c, s.apis, customBoards)

	s.engine.GET("healthz", handler.Healthz)
	s.api.GET("config", handler.Config)

	cluster := s.api.Group(":cluster")

	cluster.GET("targets", handler.ListTargets)
	cluster.GET("total-results", handler.ListTotalResults)

	cluster.GET("custom-board/:id", handler.GetCustomBoard)
	cluster.GET("custom-board/:id/cluster-resource-results", handler.ListCustomBoardClusterResourceResults)
	cluster.GET("custom-board/:id/resource-results", handler.ListCustomBoardResourceResults)
	cluster.GET("custom-board/:id/cluster-results", handler.ListCustomBoardClusterScopedResults)
	cluster.GET("custom-board/:id/results", handler.ListCustomBoardNamespaceScopedResults)

	cluster.GET("resource/:id", handler.GetResourceDetails)
	cluster.POST("resource/:id/exception", handler.CreateException)
	cluster.GET("resource/:id/results", handler.ListResourceResults)
	cluster.GET("resource/:id/resource-results", handler.ListResourceResourceResults)
	cluster.GET("results-without-resource", handler.ListResultsWithoutResource)

	ns := cluster.Group("namespace-scoped")
	ns.GET("results", handler.ListNamespaceScopedResults)
	ns.GET("resource-results", handler.ListNamespaceScopedResourceResults)

	cs := cluster.Group("cluster-scoped")
	cs.GET("results", handler.ListClusterScopedResults)
	cs.GET("resource-results", handler.ListClusterScopedResourceResults)

	cluster.GET("policy-sources", handler.ListPolicySources)
	cluster.GET("namespaces", handler.ListNamespaces)
	cluster.GET("namespace", handler.GetNamespace)
	cluster.GET(":source/policy/details", handler.GetPolicyDetails)
	cluster.GET(":source/policies", handler.ListPolicies)
	cluster.GET(":source/policy-report", handler.GetPolicyReport)
	cluster.GET(":source/namespace-report", handler.GetNamespaceReport)
	cluster.GET("dashboard", handler.GetDashboard)
	s.api.GET("clusters", handler.GetClustersDashboard)

	s.api.GET("config/:cluster/layout", handler.Layout)
}

func NewServer(engine *gin.Engine, port int, middleware []gin.HandlerFunc) *Server {
	return &Server{
		middelware: middleware,
		apis:       make(map[string]*model.Endpoints),
		engine:     engine,
		api:        engine.Group("/api", append(middleware, gzip.Gzip(gzip.DefaultCompression))...),
		proxies:    engine.Group("/proxy", middleware...),
		port:       port,
	}
}
