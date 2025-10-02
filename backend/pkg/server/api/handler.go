package api

import (
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/kyverno/policy-reporter-ui/pkg/api/model"
	"github.com/kyverno/policy-reporter-ui/pkg/auth"
	"github.com/kyverno/policy-reporter-ui/pkg/customboard"
	"github.com/kyverno/policy-reporter-ui/pkg/reports"
	"github.com/kyverno/policy-reporter-ui/pkg/service"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

type Handler struct {
	config       *Config
	clients      map[string]*model.Endpoints
	customBoards *customboard.Collection
	service      *service.Service
	reporter     *reports.ReportGenerator
}

func (h *Handler) Healthz(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) Config(ctx *gin.Context) {
	if profile := auth.ProfileFrom(ctx); profile != nil {
		cluster := h.config.Default

		clusters := make([]Cluster, 0, len(h.config.Clusters))
		for _, cl := range h.config.Clusters {
			access := cl.Allowed(profile)
			if access {
				clusters = append(clusters, cl)
			}
			if cluster == cl.Slug && !access {
				cluster = ""
			}
		}

		if cluster == "" && len(clusters) > 0 {
			cluster = clusters[0].Slug
		}

		ctx.JSON(http.StatusOK, Config{
			Clusters:    clusters,
			Default:     cluster,
			User:        h.config.User,
			Sources:     h.config.Sources,
			Banner:      h.config.Banner,
			OAuth:       h.config.OAuth,
			DisplayMode: h.config.DisplayMode,
		})
		return
	}

	ctx.JSON(http.StatusOK, h.config)
}

func (h *Handler) ListCustomBoards(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, utils.Map(h.customBoards.Boards(), MapCustomBoard))
}

func (h *Handler) ListPolicySources(ctx *gin.Context) {
	if profile := auth.ProfileFrom(ctx); profile != nil {
		if !h.config.Boards.Allowed(profile) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

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

func (h *Handler) ListNamespaces(ctx *gin.Context) {
	namespaces, err := h.service.Namespaces(ctx, ctx.Param("cluster"), ctx.Request.URL.Query())
	if err != nil {
		zap.L().Error("failed to list namespaces", zap.String("cluster", ctx.Param("cluster")), zap.Error(err))

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, namespaces)
}

func (h *Handler) CreateException(ctx *gin.Context) {
	req := service.ExceptionRequest{
		Cluster:  ctx.Param("cluster"),
		Resource: ctx.Param("id"),
	}

	if err := ctx.BindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	details, err := h.service.CreateException(ctx, req)
	if err != nil {
		zap.L().Error(
			"failed to create exception",
			zap.String("cluster", ctx.Param("cluster")),
			zap.String("id", ctx.Param("id")),
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
	var err error

	board := h.customBoards.Board(ctx.Param("id"))
	if board == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	config := MapCustomBoard(board)

	if profile := auth.ProfileFrom(ctx); profile != nil {
		if !config.Allowed(profile) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
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

	if len(sources) == 0 {
		sources, err = endpoints.Core.ListSources(ctx, url.Values{})
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	var namespaces []string
	if len(config.Namespaces.Selector) > 0 {
		ns, err := endpoints.Core.ResolveNamespaceSelector(ctx, config.Namespaces.Selector)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		namespaces = append(config.Namespaces.List, ns...)
	} else if len(config.Namespaces.List) > 0 {
		namespaces = config.Namespaces.List
	}

	if len(namespaces) == 0 {
		ctx.JSON(http.StatusOK, service.Dashboard{
			Title:          config.Name,
			ClusterScope:   config.ClusterScope,
			Sources:        sources,
			SingleSource:   len(sources) == 1,
			MultipleSource: len(sources) > 1,
			Namespaces:     make([]string, 0),
			Display:        config.Display,
			Severities:     config.Filter.Severities,
			Status:         config.Filter.Results,
			NamespaceKinds: config.Filter.NamespaceKinds,
			ClusterKinds:   config.Filter.ClusterKinds,
			Charts: service.Charts{
				ClusterScope:   make(map[string]map[string]int),
				NamespaceScope: make(map[string]*service.ChartVariants),
				Findings:       &service.Chart{},
			},
		})
		return
	}

	query["namespaces"] = namespaces

	dashboard, err := h.service.Dashboard(ctx, service.DashboardOptions{
		Cluster:        ctx.Param("cluster"),
		Sources:        sources,
		Namespaces:     namespaces,
		Display:        config.Display,
		ClusterScope:   config.ClusterScope,
		Status:         config.Filter.Results,
		Severities:     config.Filter.Severities,
		NamespaceKinds: config.Filter.NamespaceKinds,
		ClusterKinds:   config.Filter.ClusterKinds,
	}, query)
	if err != nil {
		zap.L().Error("failed to generate dashboard", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	dashboard.FilterSources = query["sources"]
	dashboard.Title = config.Name

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

	targets, err := endpoints.Core.ListTargets(ctx)
	if err != nil {
		zap.L().Error("failed to call core API", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	boards := make(map[string]CustomBoard, 0)
	list := utils.Map(h.customBoards.Boards(), MapCustomBoard)

	if profile := auth.ProfileFrom(ctx); profile != nil {
		boards = make(map[string]CustomBoard, h.customBoards.Length())

		for _, board := range list {
			if !board.Allowed(profile) {
				continue
			}

			boards[board.ID] = board
		}

		if !h.config.Boards.Allowed(profile) {
			sources = nil
		}
	} else {
		for _, board := range list {
			boards[board.ID] = board
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"targets":      len(targets) > 0,
		"sources":      MapSourceCategoryTreeToNavi(sources),
		"policies":     MapSourcesToPolicyNavi(sources),
		"customBoards": MapCustomBoardsToNavi(boards),
		"profile":      profile,
	})
}

func (h *Handler) Dashboard(ctx *gin.Context) {
	if profile := auth.ProfileFrom(ctx); profile != nil {
		if !h.config.Boards.Allowed(profile) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	endpoints, ok := h.clients[ctx.Param("cluster")]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	g := &errgroup.Group{}

	sources, _ := ctx.GetQueryArray("sources")
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
		namespaces, err = h.service.Namespaces(ctx, ctx.Param("cluster"), url.Values{
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

	dashboard, err := h.service.Dashboard(ctx, service.DashboardOptions{
		Cluster:      ctx.Param("cluster"),
		Sources:      sources,
		Namespaces:   namespaces,
		ClusterScope: true,
	}, ctx.Request.URL.Query())
	if err != nil {
		zap.L().Error("failed to generate dashboard", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, dashboard)
}

func (h *Handler) Policies(ctx *gin.Context) {
	if profile := auth.ProfileFrom(ctx); profile != nil {
		if !h.config.Boards.Allowed(profile) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

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
			ctx.JSON(http.StatusOK, MapPluginPolicies(source, policies, list))
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

func (h *Handler) PolicyReport(ctx *gin.Context) {
	data, err := h.reporter.GeneratePerPolicy(ctx, ctx.Param("cluster"), ctx.Param("source"), reports.Filter{
		Namespaces:   ctx.Request.URL.Query()["namespaces"],
		Policies:     ctx.Request.URL.Query()["policies"],
		ClusterScope: ctx.Request.URL.Query().Get("clusterScope") != "0",
		Categories:   ctx.Request.URL.Query()["categories"],
		Kinds:        ctx.Request.URL.Query()["kinds"],
	})
	if err != nil {
		zap.L().Error("failed to load generate report data", zap.String("cluster", ctx.Param("cluster")), zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("policy-report-details.html").Funcs(funcMap).ParseFiles(path.Join(os.Getenv("KO_DATA_PATH"), "templates", "reports", "policy-report-details.html"), path.Join(os.Getenv("KO_DATA_PATH"), "templates", "reports", "mui.css"))
	if err != nil {
		zap.L().Error("failed to create template", zap.String("cluster", ctx.Param("cluster")), zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(ctx.Writer, data); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *Handler) NamespaceReport(ctx *gin.Context) {
	data, err := h.reporter.GeneratePerNamespace(ctx, ctx.Param("cluster"), ctx.Param("source"), reports.Filter{
		Namespaces: ctx.Request.URL.Query()["namespaces"],
		Policies:   ctx.Request.URL.Query()["policies"],
		Categories: ctx.Request.URL.Query()["categories"],
		Kinds:      ctx.Request.URL.Query()["kinds"],
	})
	if err != nil {
		zap.L().Error("failed to load generate report data", zap.String("cluster", ctx.Param("cluster")), zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("namespace-report-details.html").Funcs(funcMap).ParseFiles(path.Join(os.Getenv("KO_DATA_PATH"), "templates", "reports", "namespace-report-details.html"), path.Join(os.Getenv("KO_DATA_PATH"), "templates", "reports", "mui.css"))
	if err != nil {
		zap.L().Error("failed to create template", zap.String("cluster", ctx.Param("cluster")), zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(ctx.Writer, data); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func NewHandler(config *Config, apis map[string]*model.Endpoints, customBoards *customboard.Collection) *Handler {
	sources := make(map[string]model.SourceConfig, len(config.Sources))
	for _, s := range config.Sources {
		sources[s.Name] = model.SourceConfig{Results: s.Excludes.Results, Exceptions: s.Exceptions, ViewType: s.ViewType, Severities: s.Excludes.Severities}
	}

	return &Handler{config, apis, customBoards, service.New(apis, sources), reports.New(apis)}
}

var funcMap = template.FuncMap{
	"add": func(i, j int) int {
		return i + j
	},
}
