package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"

	"github.com/kyverno/policy-reporter-ui/pkg/server/api"
)

type APIHandler interface {
	Register(*gin.RouterGroup)
}

type Server struct {
	engine  *gin.Engine
	api     *gin.RouterGroup
	proxies *gin.RouterGroup
	port    int
}

func (s *Server) Start() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.port))
}

func (s *Server) RegisterUI(path string) {
	fileServer := http.FileServer(http.Dir(path))

	s.engine.NoRoute(func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func (s *Server) RegisterCluster(name string, proxies map[string]*httputil.ReverseProxy) {
	group := s.proxies.Group(slug.Make(name))

	for p, rp := range proxies {
		group.Group(p).Any("/*proxy", func(ctx *gin.Context) {
			fmt.Println("handled:" + ctx.Param("proxy"))
			req := ctx.Request.Clone(ctx)
			req.URL.Path = ctx.Param("proxy")

			rp.ServeHTTP(ctx.Writer, req)
		})
	}
}

func (s *Server) RegisterAPI(c api.Config) {
	s.api.GET("config", api.ConfigHandler(c))
}

func NewServer(engine *gin.Engine, port int) *Server {
	return &Server{
		engine:  engine,
		api:     engine.Group("/api"),
		proxies: engine.Group("/proxy"),
		port:    port,
	}
}
