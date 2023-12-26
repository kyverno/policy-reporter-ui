package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/policy-reporter-ui/pkg/core/client"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

type PolicyReports struct {
	Selector map[string]string
}

type Namespaces struct {
	Selector map[string]string
	List     []string
}

type Sources struct {
	List []string
}

type CustomBoard struct {
	Name          string        `json:"name"`
	ID            string        `json:"id"`
	Namespaces    Namespaces    `json:"-"`
	Sources       Sources       `json:"-"`
	PolicyReports PolicyReports `json:"-"`
}

type CustomBoardHandler struct {
	clients map[string]*client.Client
	configs map[string]CustomBoard
}

func (h *CustomBoardHandler) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, utils.ToList(h.configs))
}

func (h *CustomBoardHandler) Details(ctx *gin.Context) {
	config, ok := h.configs[ctx.Param("id")]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	client, ok := h.clients[ctx.Param("cluster")]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	var namespaces []string
	if len(config.Namespaces.Selector) > 0 {
		ns, err := client.ResolveNamespaceSelector(ctx, config.Namespaces.Selector)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		namespaces = append(config.Namespaces.List, ns...)
	} else if len(config.Namespaces.List) > 0 {
		namespaces = config.Namespaces.List
	}

	ctx.JSON(http.StatusOK, gin.H{
		"name":       config.Name,
		"namespaces": namespaces,
		"sources":    config.Sources.List,
		"labels":     config.PolicyReports.Selector,
	})
}

func NewCustomBoardHandler(clients map[string]*client.Client, configs map[string]CustomBoard) *CustomBoardHandler {
	return &CustomBoardHandler{clients, configs}
}
