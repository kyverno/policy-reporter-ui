package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.uber.org/zap"

	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/server/api"
)

type APIHandler interface {
	Register(*gin.RouterGroup)
}

type Server struct {
	middelware []gin.HandlerFunc
	apis       map[string]*api.Endpoints
	engine     *gin.Engine
	api        *gin.RouterGroup
	proxies    *gin.RouterGroup
	port       int
}

func (s *Server) Start() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.port))
}

func (s *Server) RegisterUI(path string) {
	fileServer := http.FileServer(http.Dir(path))

	handler := append(s.middelware, func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	s.engine.NoRoute(handler...)
}

func (s *Server) RegisterCluster(name string, client *core.Client, plugins map[string]*plugin.Client, proxy *httputil.ReverseProxy) {
	id := slug.Make(name)

	s.apis[id] = &api.Endpoints{Core: client, Plugins: plugins}
	group := s.proxies.Group(id)

	group.Group("core").Any("/*proxy", func(ctx *gin.Context) {
		req := ctx.Request.Clone(ctx)
		req.URL.Path = ctx.Param("proxy")

		proxy.ServeHTTP(ctx.Writer, req)
	})

	zap.L().Debug("cluster registered", zap.String("name", name), zap.String("id", id))
}

func (s *Server) RegisterAPI(c *api.Config, configs map[string]api.CustomBoard) {
	handler := api.NewHandler(c, s.apis, configs)

	s.api.GET("config", handler.Config)
	s.api.GET("custom-board/list", handler.ListCustomBoards)
	s.api.GET("config/:cluster/custom-board/:id", handler.GetCustomBoard)
	s.api.GET("config/:cluster/resource/:id", handler.GetResourceDetails)
	s.api.GET("config/:cluster/policy-sources", handler.ListPolicySources)
	s.api.GET("config/:cluster/:source/policy/details", handler.GetPolicyDetails)
	s.api.GET("config/:cluster/:source/policies", handler.Policies)

	s.api.GET("config/:cluster/layout", handler.Layout)
	s.api.GET("config/:cluster/dashboard", handler.Dashboard)
}

func NewServer(engine *gin.Engine, port int, middleware []gin.HandlerFunc) *Server {
	return &Server{
		middelware: middleware,
		apis:       make(map[string]*api.Endpoints),
		engine:     engine,
		api:        engine.Group("/api", middleware...),
		proxies:    engine.Group("/proxy", middleware...),
		port:       port,
	}
}
