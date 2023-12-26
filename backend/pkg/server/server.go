package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.uber.org/zap"

	"github.com/kyverno/policy-reporter-ui/pkg/core/client"
	"github.com/kyverno/policy-reporter-ui/pkg/server/api"
)

type APIHandler interface {
	Register(*gin.RouterGroup)
}

type Server struct {
	clients map[string]*client.Client
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

func (s *Server) RegisterCluster(name string, client *client.Client, proxies map[string]*httputil.ReverseProxy) {
	id := slug.Make(name)

	s.clients[id] = client
	group := s.proxies.Group(id)

	for p, rp := range proxies {
		group.Group(p).Any("/*proxy", func(ctx *gin.Context) {
			fmt.Println("handled:" + ctx.Param("proxy"))
			req := ctx.Request.Clone(ctx)
			req.URL.Path = ctx.Param("proxy")

			rp.ServeHTTP(ctx.Writer, req)
		})
	}

	zap.L().Debug("cluster registered", zap.String("name", name), zap.String("id", id))
}

func (s *Server) RegisterCustomBoards(configs map[string]api.CustomBoard) {
	handler := api.NewCustomBoardHandler(s.clients, configs)

	s.api.GET("custom-board/list", handler.List)
	s.api.GET("custom-board/:cluster/:id", handler.Details)
}

func (s *Server) RegisterAPI(c api.Config) {
	s.api.GET("config", api.ConfigHandler(c))
}

func NewServer(engine *gin.Engine, port int) *Server {
	return &Server{
		clients: make(map[string]*client.Client),
		engine:  engine,
		api:     engine.Group("/api"),
		proxies: engine.Group("/proxy"),
		port:    port,
	}
}
