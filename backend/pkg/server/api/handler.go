package api

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/service"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Endpoints struct {
	Core    *core.Client
	Plugins map[string]*plugin.Client
}

type Handler struct {
	config  *Config
	clients map[string]*Endpoints
	boards  map[string]CustomBoard
	service *service.Service
}

func (h *Handler) Config(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, h.config)
}

func (h *Handler) ListCustomBoards(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, utils.ToList(h.boards))
}

func (h *Handler) ListPolicySources(ctx *gin.Context) {
	details, err := h.service.PolicySources(ctx, ctx.Param("cluster"), ctx.Request.URL.Query())
	if err != nil {
		zap.L().Error(
			"failed to generate policy sources",
			zap.String("cluster", ctx.Param("cluster")),
			zap.Error(err),
		)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, details)
}

func (h *Handler) GetResourceDetails(ctx *gin.Context) {
	details, err := h.service.ResourceDetails(ctx, ctx.Param("cluster"), ctx.Param("id"), ctx.Request.URL.Query())
	if err != nil {
		zap.L().Error(
			"failed to generate resource details",
			zap.String("cluster", ctx.Param("cluster")),
			zap.String("id", ctx.Param("id")),
			zap.Error(err),
		)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, details)
}

func (h *Handler) GetCustomBoard(ctx *gin.Context) {
	config, ok := h.boards[ctx.Param("id")]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	endpoints, ok := h.clients[ctx.Param("cluster")]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	query := ctx.Request.URL.Query()

	sources := config.Sources.List
	if len(sources) > 0 {
		query["sources"] = sources
	}

	g := &errgroup.Group{}
	if len(sources) == 0 {
		g.Go(func() error {
			var err error
			sources, err = endpoints.Core.ListSources(ctx, url.Values{})

			return err
		})
	}

	var namespaces []string
	if len(config.Namespaces.Selector) > 0 {
		ns, err := endpoints.Core.ResolveNamespaceSelector(ctx, config.Namespaces.Selector)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		namespaces = append(config.Namespaces.List, ns...)
	} else if len(config.Namespaces.List) > 0 {
		namespaces = config.Namespaces.List
	}

	if len(namespaces) > 0 {
		query["namespaces"] = namespaces
	}

	dashboard, err := h.service.Dashboard(ctx, ctx.Param("cluster"), sources, namespaces, config.ClusterScope, query)
	if err != nil {
		zap.L().Error("failed to generate dashboard", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	dashboard.FilterSources = query["sources"]

	ctx.JSON(http.StatusOK, dashboard)
}

func (h *Handler) Layout(ctx *gin.Context) {
	endpoints, ok := h.clients[ctx.Param("cluster")]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	sources, err := endpoints.Core.ListSourceCategoryTree(ctx, ctx.Request.URL.Query())
	if err != nil {
		zap.L().Error("failed to call core API", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	profile, _ := ctx.Get("profile")

	ctx.JSON(http.StatusOK, gin.H{
		"sources":      MapSourceCategoryTreeToNavi(sources),
		"policies":     MapSourcesToPolicyNavi(sources),
		"customBoards": MapCustomBoardsToNavi(h.boards),
		"profile":      profile,
	})
}

func (h *Handler) Dashboard(ctx *gin.Context) {
	endpoints, ok := h.clients[ctx.Param("cluster")]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	g := &errgroup.Group{}

	sources, ok := ctx.GetQueryArray("sources")
	if len(sources) == 0 {
		g.Go(func() error {
			var err error
			sources, err = endpoints.Core.ListSources(ctx, url.Values{})

			return err
		})
	}

	query := ctx.Request.URL.Query()

	var namespaces []string
	g.Go(func() error {
		var err error
		namespaces, err = endpoints.Core.ListNamespaces(ctx, url.Values{
			"sources":    query["sources"],
			"kinds":      query["kinds"],
			"categories": query["categories"],
			"policies":   query["policies"],
		})

		return err
	})

	if err := g.Wait(); err != nil {
		zap.L().Error("failed to call core api", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	dashboard, err := h.service.Dashboard(ctx, ctx.Param("cluster"), sources, namespaces, true, ctx.Request.URL.Query())
	if err != nil {
		zap.L().Error("failed to generate dashboard", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, dashboard)
}

func (h *Handler) Policies(ctx *gin.Context) {
	endpoints, ok := h.clients[ctx.Param("cluster")]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	source := ctx.Param("source")

	query := ctx.Request.URL.Query()
	query.Set("sources", source)

	list, err := endpoints.Core.ListPolicies(ctx, query)
	if err != nil {
		zap.L().Error("failed to load policies from core api", zap.String("cluster", ctx.Param("cluster")), zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if plugin, ok := endpoints.Plugins[source]; ok {
		policies, err := plugin.GetPolicies(ctx)
		if err != nil {
			zap.L().Error("failed to load policies from plugin", zap.String("cluster", ctx.Param("cluster")), zap.String("plugin", source), zap.Error(err))
		} else {
			ctx.JSON(http.StatusOK, MapPluginPolicies(policies, list))
			return
		}
	}

	ctx.JSON(http.StatusOK, MapPoliciesFromCore(list))
}

func (h *Handler) GetPolicyDetails(ctx *gin.Context) {
	details, err := h.service.PolicyDetails(ctx, ctx.Param("cluster"), ctx.Param("source"), ctx.Query("policies"), ctx.Request.URL.Query())
	if err != nil {
		zap.L().Error(
			"failed to generate policy sources",
			zap.String("cluster", ctx.Param("cluster")),
			zap.Error(err),
		)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, details)
}

func NewHandler(config *Config, apis map[string]*Endpoints, customBoards map[string]CustomBoard) *Handler {
	endpoints := make(map[string]*service.Endpoints, len(apis))
	for cluster, value := range apis {
		endpoints[cluster] = &service.Endpoints{
			Core:    value.Core,
			Plugins: value.Plugins,
		}
	}

	return &Handler{config, apis, customBoards, service.New(endpoints)}
}
